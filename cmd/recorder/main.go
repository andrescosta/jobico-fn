package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/recorder/service"
)

func main() {
	err := service.Service{}.Start(context.Background())
	if err != nil {
		log.Panicf("error starting recorder service: %s", err)
	}
}
