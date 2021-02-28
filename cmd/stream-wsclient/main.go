package main

import (
	"bytes"
	"context"
	"fmt"

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

func main() {
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
