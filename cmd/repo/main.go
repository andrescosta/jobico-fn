package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/repo/service"
)

func main() {
	err := service.Service{}.Start(context.Background())
	if err != nil {
		log.Panicf("error starting repo service: %s", err)
	}
}
