package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/executor/service"
)

func main() {
	err := service.Service{}.Start(context.Background())
	if err != nil {
		log.Fatalf("error starting executors: %v", err)
	}
}
