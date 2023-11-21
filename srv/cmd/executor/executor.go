package main

import (
	"context"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/workflew/srv/internal/executor"
	"github.com/andrescosta/workflew/srv/internal/queue"
)

func main() {
	service.Start(serviceFunc)
}

func serviceFunc(ctx context.Context) error {
	dirq := "../queue/" + queue.DIR + "/"

	if err := executor.StartExecutors(ctx, dirq); err != nil {
		return err
	}
	return nil
}
