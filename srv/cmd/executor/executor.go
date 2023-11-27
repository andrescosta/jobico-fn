package main

import (
	"context"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/workflew/srv/internal/executor"
)

func main() {
	service.StartNamed("executor", serviceFunc)
}

func serviceFunc(ctx context.Context) error {
	if err := executor.StartExecutors(ctx); err != nil {
		return err
	}
	return nil
}
