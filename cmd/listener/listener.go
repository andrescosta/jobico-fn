package main

import (
	"fmt"
	"log"

	"github.com/andrescosta/goico/pkg/context"
	"github.com/andrescosta/jobico/cmd/listener/service"
)

func main() {
	ctx, cancel := context.ForEndSignals()
	defer cancel()
	svc, err := service.New(ctx)
	if err != nil {
		log.Panicf("error creating listener service: %v", err)
	}
	defer func() {
		if err := svc.Dispose(); err != nil {
			fmt.Printf("error disposing listener resources %v", err)
		}
	}()
	if err := svc.Start(); err != nil {
		log.Panicf("error starting listener service: %v", err)
	}
}
