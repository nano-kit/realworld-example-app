package main

import (
	"context"
	"fmt"
	"realworld-example-app/internal/json"
	realworld "realworld-example-app/proto/realworld"

	gcli "github.com/micro/go-micro/v2/client/grpc"
)

func main() {
	client := gcli.NewClient()
	service := realworld.NewRealworldService("com.example.service.realworld", client)
	ctx := context.Background()
	stream, err := service.PingPong(ctx)
	if err != nil {
		panic(err)
	}
	limit := int64(5)
	stroke := int64(1)
	for {
		if err := stream.Send(&realworld.Ping{Stroke: stroke}); err != nil {
			panic(err)
		}
		pong, err := stream.Recv()
		if err != nil {
			panic(err)
		}
		fmt.Println(json.Stringify(pong))

		stroke++
		if stroke > limit {
			break
		}
	}
	stream.Close()
}
