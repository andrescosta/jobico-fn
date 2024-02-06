package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/executor/service"
)

func main() {
	svc, err := service.New(context.Background())
	if err != nil {
		log.Panicf("error creating executor service: %s", err)
	}
	if err := svc.Start(); err != nil {
		log.Panicf("error starting executor service: %s", err)
	}
}
