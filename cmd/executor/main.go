package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andrescosta/jobico/cmd/executor/service"
)

func main() {
	svc, err := service.New(context.Background())
	if err != nil {
		log.Panicf("error creating executor service: %s", err)
	}
	defer func() {
		if err := svc.Dispose(); err != nil {
			fmt.Printf("error disposing resources %v", err)
		}
	}()
	if err := svc.Start(); err != nil {
		log.Panicf("error starting executor service: %s", err)
	}
}
