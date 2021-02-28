package handler

import (
	"context"
	"fmt"
	"math"
	"realworld-example-app/internal/json"
	realworld "realworld-example-app/proto/realworld"

	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
)

type Realworld struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Realworld) Call(ctx context.Context, req *realworld.Request, rsp *realworld.Response) error {
	log.Info("Received Realworld.Call request", json.Stringify(req))

	if req.Age < 0 {
		// return a raw go error
		// client should get {
		//  	"id":"go.micro.client",
		//		"code":500,
		//		"detail":"invalid age: -1",
		// 		"status":"Internal Server Error"
		// 	}
		return fmt.Errorf("invalid age: %v", req.Age)
	}

	if req.Age == 0 {
		// return a go-micro rpc aware custom error
		// client should get {
		// 		"id":"zero-age",
		//		"code":400,
		//		"detail":"age is 0, forget to set the age?",
		//		"status":"Bad Request"
		//	}
		return errors.BadRequest("zero-age", "age is %v, forget to set the age?", req.Age)
	}

	if req.Age > 200 {
		return errors.New("age-too-old", fmt.Sprintf("age is %v, too old!", req.Age), 520)
	}

	rsp.Msg = "Hello " + req.Name
	rsp.NumInt32 = math.MaxInt32
	rsp.NumInt64 = math.MaxInt64
	rsp.NumFloat = math.MaxFloat32
	rsp.NumDouble = math.MaxFloat64
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Realworld) Stream(ctx context.Context, req *realworld.StreamingRequest, stream realworld.Realworld_StreamStream) error {
	log.Infof("Received Realworld.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&realworld.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Realworld) PingPong(ctx context.Context, stream realworld.Realworld_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&realworld.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
