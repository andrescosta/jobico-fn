package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/server"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/srv/internal/queue"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const Name = "Queue"

func main() {
	service.StartNamed("Queue", serviceFunc)
}

func serviceFunc(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	s := grpc.NewServer()
	pb.RegisterQueueServer(s, &queue.Server{})
	reflection.Register(s)

	srv, err := server.New(os.Getenv("queue.addr"))
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}
	logger.Info().Msgf("Started at:%s", srv.Addr())
	err = srv.ServeGRPC(ctx, s)
	logger.Info().Msg("Stopped")
	return err
}
