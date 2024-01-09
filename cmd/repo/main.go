package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/repo/server"
)

func main() {
	svc, err := grpc.New(
		grpc.WithName("repo"),
		grpc.WithContext(context.Background()),
		grpc.WithServiceDesc(&pb.Repo_ServiceDesc),
		grpc.WithInitHandler(func(ctx context.Context) (any, error) {
			return server.New(ctx, env.String("repo.dir", "./")), nil
		}),
	)
	if err != nil {
		log.Panicf("error starting repo service: %s", err)
	}

	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving repo service: %s", err)
	}
}
