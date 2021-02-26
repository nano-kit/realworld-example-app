package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	realworld "realworld-example-app/proto/realworld"
)

type Realworld struct{}

func (e *Realworld) Handle(ctx context.Context, msg *realworld.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *realworld.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
