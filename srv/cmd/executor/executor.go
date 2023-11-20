package main

import (
	"context"

	"github.com/andrescosta/workflew/srv/internal/executor"
	"github.com/andrescosta/workflew/srv/internal/queue"
	"github.com/andrescosta/workflew/srv/pkg/service"
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
