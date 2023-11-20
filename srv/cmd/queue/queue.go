package main

import (
	"context"
	"fmt"
	"os"

	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/srv/internal/queue"
	"github.com/andrescosta/workflew/srv/pkg/server"
	"github.com/andrescosta/workflew/srv/pkg/service"
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
	pb.RegisterQueueServer(s, &queue.Server{})
	reflection.Register(s)

	srv, err := server.New(os.Getenv("queue.port"))
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}
	logger.Info().Msgf("Queue started at:%s", srv.Addr())
	err = srv.ServeGRPC(ctx, s)
	logger.Info().Msg("Queue stopped")
	return err
}
