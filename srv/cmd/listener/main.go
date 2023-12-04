package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/jobico/srv/internal/listener"
	"github.com/go-chi/chi/v5"
)

func main() {
	service.NewHttpService(context.Background(), "listener", "events",
		func(ctx context.Context) chi.Router {
			controler, err := listener.New(ctx)
			if err != nil {
				log.Fatal(err)
			}
			return controler.Routes(ctx)
		}).Serve()
}
