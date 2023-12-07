package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service"

	//"github.com/andrescosta/goico/pkg/service/obs"
	"github.com/andrescosta/jobico/srv/internal/listener"
)

func main() {
	svc, err := service.NewHttpService(context.Background(), "listener", listener.ConfigureRoutes)
	if err != nil {
		log.Fatal(err)
	}

	if err := svc.Serve(); err != nil {
		log.Fatal(err)
	}
}
