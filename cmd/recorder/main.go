package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/recorder/service"
)

func main() {
	svc, err := service.New(context.Background())
	if err != nil {
		log.Panicf("error creating recorder service: %s", err)
	}
	if err := svc.Start(); err != nil {
		log.Panicf("error starting recorder service: %s", err)
	}
}
