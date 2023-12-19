package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/queue"
)

const Name = "Queue"

func main() {
	svc, err := grpc.New(
		grpc.WithName("queue"),
		grpc.WithContext(context.Background()),
		grpc.WithServiceDesc(&pb.Queue_ServiceDesc),
		grpc.WithInitHandler(func(ctx context.Context) (any, error) {
			return queue.NewServer(ctx)
		}),
	)
	if err != nil {
		log.Panicf("error starting queue service because: %s", err)
	}
	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving queue service  because: %s", err)
	}
}
