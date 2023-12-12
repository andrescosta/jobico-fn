package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/jobico/srv/internal/executor"
)

func main() {
	e, err := service.NewHeadlessService(context.Background(), "executor",
		func(ctx context.Context) error {
			m, err := executor.NewExecutorMachine(ctx)
			if err != nil {
				return err
			}
			if err := m.StartExecutors(ctx); err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Start(); err != nil {
		log.Fatalf("error running the executor service %s", err)
	}
}
