package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/queue/service"
)

const Name = "Queue"

func main() {
	err := service.Service{}.Start(context.Background())
	if err != nil {
		log.Fatalf("error starting queue service because: %v", err)
	}
}
