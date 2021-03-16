package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./html")))

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
