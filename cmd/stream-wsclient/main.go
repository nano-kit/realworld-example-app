package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"realworld-example-app/internal/json"
	protocol "realworld-example-app/proto/realworld"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/protobuf/jsonpb"
	pb "github.com/golang/protobuf/proto"
)

var (
	jsonMarshaler = jsonpb.Marshaler{
		OrigName: true,
	}
	jsonUnmarshaler = jsonpb.Unmarshaler{}
)

func jsonMarshal(m pb.Message) ([]byte, error) {
	b := new(bytes.Buffer)
	err := jsonMarshaler.Marshal(b, m)
	return b.Bytes(), err
}

func jsonUnmarshal(data []byte, m pb.Message) error {
	return jsonUnmarshaler.Unmarshal(bytes.NewReader(data), m)
}

func main1() {
	url := "ws://127.0.0.1:8080/realworld/Realworld/PingPong"
	dialer := ws.Dialer{}
	ctx := context.Background()
	conn, _, _, err := dialer.Dial(ctx, url)
	if err != nil {
		panic(err)
	}
	limit := int64(5)
	stroke := int64(1)
	for {
		if err := wsutil.WriteClientText(conn, []byte(fmt.Sprintf(`{"stroke":%d}`, stroke))); err != nil {
			panic(err)
		}
		buf, err := wsutil.ReadServerText(conn)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(buf))

		stroke++
		if stroke > limit {
			break
		}
	}
	conn.Close()
}

func main() {
	url := "ws://127.0.0.1:8080/realworld/Clubhouse/Subscribe"
	dialer := ws.Dialer{}
	ctx := context.Background()
	conn, _, _, err := dialer.Dial(ctx, url)
	if err != nil {
		panic(err)
	}

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsInNjb3BlcyI6WyJub3JtYWwiXSwibWV0YWRhdGEiOm51bGwsImV4cCI6MTYxNTcxNDMyMCwiaXNzIjoiY29tLmV4YW1wbGUiLCJzdWIiOiJ1c2VyMDAxIn0.dz7qchVXY2wOzIDIgja8okGh8jrrRyQ14b6q73rD9S-EV4mgo1Kc0BMaZtTnqJoD0G2Na3X46y6-murPSiu1n1OcMuUuYtucv8T9CzCvYzNwqU3MvdJNEBN56jCJW8KfxZhiK0r9_UYRWeli5ysdHWApmuOKAd1Hu4Fp9Ompwxitxn71FJcQy9TB9RItVEdL6JZyrprOX4W_A7YWb6c4eSTo_SIaTa3yi0ynyXsjrC-yWyPeudhgd6Jgkbjwcs-Q3Jt1PIXVD_Xgx9BbG0hJxxYcYo0pqPq5_HvTsIGf0Pw3euHBzJ-bYweu-Ln1yVS5qP2shb9d_L8vAfuHmsbIEA"

	if err := wsutil.WriteClientText(conn, []byte(fmt.Sprintf(`{"token":"%s"}`, token))); err != nil {
		panic(err)
	}

	for {
		buf, err := wsutil.ReadServerText(conn)
		if err != nil {
			panic(err)
		}
		log.Println(string(buf))

		var push protocol.ServerPush
		if err := jsonUnmarshal(buf, &push); err != nil {
			panic(err)
		}
		if push.T == protocol.ServerPush_HEARTBEAT {
			if err := wsutil.WriteClientText(conn, []byte("{}")); err != nil {
				panic(err)
			}
			continue
		}
		log.Printf("%v", json.Stringify(push))
	}
}
