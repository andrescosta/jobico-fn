package main

import (
	"fmt"
	"log"

	"github.com/andrescosta/goico/pkg/context"
	"github.com/andrescosta/jobico/cmd/executor/service"
)

func main() {
	ctx, cancel := context.ForEndSignals()
	defer cancel()
	svc, err := service.New(ctx)
	if err != nil {
		log.Panicf("error creating executor service: %s", err)
	}
	defer func() {
		if err := svc.Dispose(); err != nil {
			fmt.Printf("error disposing executor resources %v", err)
		}
	}()
	if err := svc.Start(); err != nil {
		log.Panicf("error starting executor service: %s", err)
	}
}
