package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/srv/internal/queue"
)

const Name = "Queue"

func main() {
	svc, err := service.NewGrpService(context.Background(), "queue",
		&pb.Queue_ServiceDesc,
		func(ctx context.Context) (any, error) {
			return &queue.Server{}, nil
		},
		service.EmptyhealthCheckHandler)

	if err != nil {
		log.Panicf("error starting queue service %s", err)
	}
	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving queue service %s", err)
	}
}
