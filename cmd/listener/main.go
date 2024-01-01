package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service/http"
	"github.com/andrescosta/jobico/internal/listener"
)

func main() {
	svc, err := http.New(
		http.WithContext(context.Background()),
		http.WithName("listener"),
		http.WithInitRoutesFn(listener.ConfigureRoutes),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := svc.Serve(); err != nil {
		log.Fatal(err)
	}
}
