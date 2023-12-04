package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/api/types"
	repo "github.com/andrescosta/jobico/repo/internal"
)

func main() {
	svc, err := service.NewGrpcService(context.Background(), "repo",
		&pb.Repo_ServiceDesc, func(ctx context.Context) (any, error) {
			return repo.NewServer(ctx, env.GetAsString("repo.dir", "./")), nil
		}, nil)

	if err != nil {
		log.Panicf("error starting repo service: %s", err)
	}
	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving repo service: %s", err)
	}

}
