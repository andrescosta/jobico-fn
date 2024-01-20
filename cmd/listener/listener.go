package main

import (
	"context"
	"log"

	"github.com/andrescosta/jobico/cmd/listener/service"
)

func main() {
	err := service.Service{}.Start(context.Background())
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
