package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/recorder/internal/recorder"
)

func main() {
	svc, err := grpc.New(
		grpc.WithName("recorder"),
		grpc.WithContext(context.Background()),
		grpc.WithServiceDesc(&pb.Recorder_ServiceDesc),
		grpc.WithInitHandler(func(ctx context.Context) (any, error) {
			return recorder.NewServer(".\\log.log")
		}),
	)
	if err != nil {
		log.Panicf("error starting recorder service: %s", err)
	}
	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving recorder service: %s", err)
	}
}
