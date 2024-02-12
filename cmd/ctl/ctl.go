package main

import (
	"log"

	"github.com/andrescosta/goico/pkg/context"
	"github.com/andrescosta/jobico/cmd/ctl/service"
)

func main() {
	ctx, cancel := context.ForEndSignals()
	defer cancel()
	svc, err := service.New(ctx)
	if err != nil {
		log.Panicf("error creating ctl service: %s", err)
	}
	defer svc.Dispose()
	if err := svc.Start(); err != nil {
		log.Panicf("error starting ctl service: %s", err)
	}
}
