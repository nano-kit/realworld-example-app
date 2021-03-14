package main

import (
	"realworld-example-app/handler"
	"realworld-example-app/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	realworld "realworld-example-app/proto/realworld"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.example.service.realworld"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	realworld.RegisterRealworldHandler(service.Server(), new(handler.Realworld))
	realworld.RegisterClubhouseHandler(service.Server(), handler.NewClubhouse())

	// Register Struct as Subscriber
	micro.RegisterSubscriber("com.example.service.realworld", service.Server(), new(subscriber.Realworld))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
