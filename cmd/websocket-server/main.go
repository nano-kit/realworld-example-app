package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/store/sqlite"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	// init http client
	hc := &http.Client{Timeout: 10 * time.Second}

	// init go micro client
	cc := &clientWrapper{grpc.NewClient(), "com.example"}

	// init oauth for github
	oauthGithub := new(oauthGithub)
	oauthGithub.Init(hc, cc)

	// init token refresher
	tokenRef := new(tokenRefresher)
	tokenRef.Init(cc)

	// init url router
	r := mux.NewRouter()

	// setup route of oauth github
	oauthGithubRouter := r.PathPrefix("/oauth/github").Subrouter()
	oauthGithubRouter.
		HandleFunc("/login", oauthGithub.Login).
		Methods(http.MethodGet)
	oauthGithubRouter.
		HandleFunc("/callback", oauthGithub.Callback).
		Methods(http.MethodGet)

	// setup route of token refresher
	r.HandleFunc("/token", tokenRef.Token)

	// setup route of html pages
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./html")))

	// start micro web service
	service := web.NewService(
		web.Name("com.example.portal"),
		web.Address(":9100"),
		web.Handler(r),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

type tokenRefresher struct {
	cc client.Client
}

func (t *tokenRefresher) Init(cc client.Client) {
	t.cc = cc
}

// Token generates access token by refresh token
func (t *tokenRefresher) Token(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		http.Error(w, errors.Unauthorized("need-login", "need login to proceed with the API").Error(), http.StatusUnauthorized)
		return
	}

	authSrv := pb.NewAuthService("go.micro.auth", t.cc)
	res, err := authSrv.Token(context.Background(), &pb.TokenRequest{
		RefreshToken: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(res.GetToken().GetAccessToken()))
}

type oauthGithub struct {
	ClientID     string
	ClientSecret string

	state store.Store
	hc    *http.Client
	cc    client.Client
}

func (o *oauthGithub) Init(hc *http.Client, cc client.Client) {
	// read client secrets from file, which can not be put to github
	conf, err := config.NewConfig(
		config.WithSource(file.NewSource(file.WithPath("oauthGithub.json"))),
	)
	if err != nil {
		panic(err)
	}
	if err := conf.Scan(o); err != nil {
		panic(err)
	}

	// init a store for state management
	o.state = sqlite.NewStore(
		store.Database("oauth"),
		store.Table("state"),
	)

	// init http client
	o.hc = hc

	// init go micro client
	o.cc = cc
}

type clientWrapper struct {
	client.Client
	namespace string
}

func (a *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if len(a.namespace) > 0 {
		ctx = metadata.Set(ctx, "Micro-Namespace", a.namespace)
	}
	return a.Client.Call(ctx, req, rsp, opts...)
}

func (o *oauthGithub) Login(w http.ResponseWriter, r *http.Request) {
	// generate an unguessable random string as state
	state := uuid.NewString()
	record := &store.Record{
		Key:    state,
		Expiry: 10 * time.Minute,
	}
	if err := o.state.Write(record); err != nil {
		log.Errorf("can not save state: %v", err)
		return
	}

	u, _ := url.Parse("https://github.com/login/oauth/authorize")
	q := u.Query()
	q.Set("client_id", o.ClientID)
	q.Set("state", state)
	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (o *oauthGithub) Callback(w http.ResponseWriter, r *http.Request) {
	// get code and state
	q := r.URL.Query()
	code := q.Get("code")
	state := q.Get("state")

	// check state
	if _, err := o.state.Read(state); err != nil {
		log.Errorf("validate state failure: %v", err)
		return
	}
	o.state.Delete(state)

	// exchange this code for an access token
	token, err := o.getAccessToken(code, state)
	if err != nil {
		log.Errorf("exchange for an access token: %v", err)
		return
	}

	// get user profile using the access token
	gu, err := o.getUserProfile(token)
	if err != nil {
		log.Errorf("get github user profile: %v", err)
		return
	}

	// create or get my user account
	acc, err := o.getUserAccount(gu)
	if err != nil {
		log.Errorf("get my user account: %v", err)
		return
	}

	// get my token by account
	mytoken, err := o.getUserToken(acc)
	if err != nil {
		log.Errorf("get my token: %v", err)
	}
	/*
		log.Infof(`Login with GitHub successfully,
		ID: %s
		Name: %s
		Email: %s
		Avatar: %s
		Location: %s
		Company: %s
		GitHub Access Token: %s
		My Access Token: %s
		My Refresh Token: %s
		My Access Token Expiry: %s
		My Password: %s`,
			acc.Id,
			acc.Metadata["Name"],
			acc.Metadata["Email"],
			acc.Metadata["AvatarUrl"],
			acc.Metadata["Location"],
			acc.Metadata["Company"],
			token,
			mytoken.AccessToken,
			mytoken.RefreshToken,
			time.Unix(mytoken.Expiry, 0).Local().Format(time.RFC3339),
			acc.Secret)
	*/

	// save token to cookie
	// https://dev.to/cotter/localstorage-vs-cookies-all-you-need-to-know-about-storing-jwt-tokens-securely-in-the-front-end-15id
	cookie := http.Cookie{
		Name:     "refresh-token",
		Value:    mytoken.RefreshToken,
		Path:     "/token",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	// redirect to the index page

}

func (o *oauthGithub) getAccessToken(code, state string) (token string, err error) {
	u := "https://github.com/login/oauth/access_token"
	data := url.Values{}
	data.Set("client_id", o.ClientID)
	data.Set("client_secret", o.ClientSecret)
	data.Set("code", code)
	data.Set("state", state)
	r, err := http.NewRequest("POST", u, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept", "application/json")
	res, err := o.hc.Do(r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var resb struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resb); err != nil {
		return "", err
	}
	return resb.AccessToken, nil
}

type GitHubUser struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	Location  string `json:"location"`
	Email     string `json:"email"`
}

func (o *oauthGithub) getUserProfile(token string) (*GitHubUser, error) {
	u := "https://api.github.com/user"
	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "token "+token)
	res, err := o.hc.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var gu GitHubUser
	if err := json.NewDecoder(res.Body).Decode(&gu); err != nil {
		return nil, err
	}
	return &gu, nil
}

func (o *oauthGithub) getUserAccount(gu *GitHubUser) (*pb.Account, error) {
	authSrv := pb.NewAuthService("go.micro.auth", o.cc)
	res, err := authSrv.Generate(context.Background(), &pb.GenerateRequest{
		Id: gu.Login,
		Metadata: map[string]string{
			"Name":      gu.Name,
			"Company":   gu.Company,
			"AvatarUrl": gu.AvatarUrl,
			"Email":     gu.Email,
			"Location":  gu.Location,
		},
		Scopes:   []string{"basic"},
		Provider: "oauth",
		Type:     "user",
		Secret:   uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}
	return res.GetAccount(), nil
}

func (o *oauthGithub) getUserToken(acc *pb.Account) (*pb.Token, error) {
	authSrv := pb.NewAuthService("go.micro.auth", o.cc)
	res, err := authSrv.Token(context.Background(), &pb.TokenRequest{
		Id:     acc.Id,
		Secret: acc.Secret,
	})
	if err != nil {
		return nil, err
	}
	return res.GetToken(), nil
}
