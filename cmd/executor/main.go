package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service/process"
	"github.com/andrescosta/jobico/internal/executor"
)

func main() {
	e, err := process.New(
		process.WithContext(context.Background()),
		process.WithName("executor"),
		process.WithServeHandler(func(ctx context.Context) error {
			m, err := executor.NewExecutorMachine(ctx)
			if err != nil {
				return err
			}
			if err := m.StartExecutors(ctx); err != nil {
				return err
			}
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Serve(); err != nil {
		log.Fatalf("error running the executor service %s", err)
	}
}
