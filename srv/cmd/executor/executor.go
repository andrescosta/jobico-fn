package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/workflew/srv/internal/executor"
)

func main() {
	if err := service.NewHeadlessService(context.Background(), "executor",
		func(ctx context.Context) error {
			if err := executor.StartExecutors(ctx); err != nil {
				return err
			}
			return nil
		}).Serve(); err != nil {
		log.Fatalf("error running executor service %s", err)
	}
}
