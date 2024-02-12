package main

import (
	"log"

	"github.com/andrescosta/goico/pkg/context"
	"github.com/andrescosta/jobico/cmd/recorder/service"
)

func main() {
	ctx, cancel := context.ForEndSignals()
	defer cancel()
	svc, err := service.New(ctx)
	if err != nil {
		log.Panicf("error creating recorder service: %s", err)
	}
	defer svc.Dispose()
	if err := svc.Start(); err != nil {
		log.Panicf("error starting recorder service: %s", err)
	}
}
