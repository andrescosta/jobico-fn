package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/workflew/api/types"
	repo "github.com/andrescosta/workflew/repo/internal"
)

func main() {
	svc, err := service.NewGrpService(context.Background(), "repo",
		&pb.Repo_ServiceDesc, func(ctx context.Context) (any, error) {
			return &repo.Server{
				Repo: &repo.FileRepo{
					Dir: env.GetAsString("repo.dir", "./"),
				},
			}, nil
		})

	if err != nil {
		log.Panicf("error starting repo service: %s", err)
	}
	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving repo service: %s", err)
	}

}
