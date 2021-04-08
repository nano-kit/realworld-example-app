package handler

import (
	"context"
	"fmt"
	"math"
	"realworld-example-app/internal/json"
	realworld "realworld-example-app/proto/realworld"
	"strings"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
)

type Realworld struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Realworld) Call(ctx context.Context, req *realworld.Request, rsp *realworld.Response) error {
	acc, err := auth.AccountFromContext(ctx)
	log.Infof("Received Realworld.Call request %v from account %+v, err=%v", json.Stringify(req), acc, err)

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

func (e *Realworld) Upload(ctx context.Context, stream realworld.Realworld_UploadStream) error {
	var file []string
	for {
		datapack, err := stream.Recv()
		if err != nil {
			return err
		}
		file = append(file, datapack.Line)
		if datapack.Done {
			break
		}
	}
	return stream.SendMsg(&realworld.UploadResp{
		TotalLines: int64(len(file)),
		File:       strings.Join(file, "\n"),
	})
}
