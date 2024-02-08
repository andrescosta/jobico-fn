package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andrescosta/jobico/cmd/listener/service"
)

func main() {
	svc, err := service.New(context.Background())
	if err != nil {
		log.Panicf("error creating listener service: %s", err)
	}
	defer func() {
		if err := svc.Dispose(); err != nil {
			fmt.Printf("error disposing resources %v", err)
		}
	}()
	if err := svc.Start(); err != nil {
		log.Panicf("error starting listener service: %s", err)
	}
}
