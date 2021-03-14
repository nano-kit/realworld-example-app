package auth

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	microauth "github.com/micro/go-micro/v2/auth"
)

// authClaims to be encoded in the JWT
type authClaims struct {
	Type     string            `json:"type"`
	Scopes   []string          `json:"scopes"`
	Metadata map[string]string `json:"metadata"`

	jwt.StandardClaims
}

// AccountFromToken restore account from the access JWT token
func AccountFromToken(token string) (*microauth.Account, bool) {
	// check token format
	if len(strings.Split(token, ".")) != 3 {
		return nil, false
	}

	// get the public key from env
	key := os.Getenv("MICRO_AUTH_PUBLIC_KEY")
	if key == "" {
		return nil, false
	}

	// decode the public key
	pub, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, false
	}

	// parse the public key
	res, err := jwt.ParseWithClaims(token, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(pub)
	})
	if err != nil {
		return nil, false
	}

	// validate the token
	if !res.Valid {
		return nil, false
	}
	claims, ok := res.Claims.(*authClaims)
	if !ok {
		return nil, false
	}

	// return the token
	return &microauth.Account{
		ID:       claims.Subject,
		Issuer:   claims.Issuer,
		Type:     claims.Type,
		Scopes:   claims.Scopes,
		Metadata: claims.Metadata,
	}, true
}
