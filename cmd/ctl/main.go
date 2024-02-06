package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/ctl/service"
)

func main() {
	svc, err := service.New(context.Background())
	if err != nil {
		log.Panicf("error creating ctl service: %s", err)
	}
	if err := svc.Start(); err != nil {
		log.Panicf("error starting ctl service: %s", err)
	}
}
