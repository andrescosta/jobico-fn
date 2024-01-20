package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/ctl/service"
)

func main() {
	err := service.Service{}.Start(context.Background())
	if err != nil {
		log.Panicf("error starting ctl service: %s", err)
	}
}
