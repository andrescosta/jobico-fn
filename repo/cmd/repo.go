package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/server"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/workflew/api/types"
	repo "github.com/andrescosta/workflew/repo/internal"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	service.Start(serviceFunc)
}

func serviceFunc(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	s := grpc.NewServer()

	pb.RegisterRepoServer(s, &repo.Server{
		Repo: &repo.FileRepo{
			Dir: env.GetAsString("repo.dir", "./"),
		},
	})
	reflection.Register(s)

	srv, err := server.New(os.Getenv("repo.port"))
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}
	logger.Info().Msgf("Repo started at:%s", srv.Addr())
	err = srv.ServeGRPC(ctx, s)
	logger.Info().Msg("Repo stopped")
	return err
}
