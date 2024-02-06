package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/listener/service"
)

func main() {
	svc, err := service.New(context.Background())
	if err != nil {
		log.Panicf("error creating listener service: %s", err)
	}
	if err := svc.Start(); err != nil {
		log.Panicf("error starting listener service: %s", err)
	}
}
