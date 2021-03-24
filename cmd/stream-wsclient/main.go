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

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsInNjb3BlcyI6WyJiYXNpYyJdLCJtZXRhZGF0YSI6eyJBdmF0YXJVcmwiOiJodHRwczovL2F2YXRhcnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tL3UvMTAzMjExOTg_dj00IiwiQ29tcGFueSI6IiIsIkVtYWlsIjoiIiwiTG9jYXRpb24iOiIiLCJOYW1lIjoiIn0sImV4cCI6MTYxNjU1NjA0OSwiaXNzIjoiY29tLmV4YW1wbGUiLCJzdWIiOiJwaWFub2h1YiJ9.ToW_Z5QBCxnklYcEOuEQyfyyglzxjNhZbwMWdSHcTsDcfYX-Ynj-qbDCFq8KwwCKVeI_8shmnlS9fEDCey1cerzmRP2BHqppBJMKbMOp4g0U0q6EBa_HaVQCB6Tq9AKuPeHEGEBuNj9Cw8cFQ98TLSa7RnMnwg1k3-3jTX0S1IOoRFrrfdtnH2ZkgKP6MdVH1IU1MIde-78C12dZ4x6XI8LB-IPU9AwBiBgNAJOIB-Y5E7yXq98GmJfYP498H7APF2gAf0AkscRFBkYRCVo3vYTNMxO0iY0ofh4wkp9HaLOO9qDU75JRHl3j5yG6FWCqB_O_YidXXYCIzxMk6xh0TQ"

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
