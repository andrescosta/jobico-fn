package main

import (
	"context"
	"log"
	"net/http"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/workflew/srv/internal/listener"
)

func main() {
	service.NewHttpService(context.Background(), "listener",
		func(ctx context.Context) http.Handler {
			controler, err := listener.New(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			return controler.Routes(ctx)
		}).Serve()
}
