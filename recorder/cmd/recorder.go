package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/server"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/recorder/internal/recorder"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	service.StartNamed("Recorder", serviceFunc)
}

func serviceFunc(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	s := grpc.NewServer()

	pb.RegisterRecorderServer(s, recorder.NewServer(ctx, ".\\log.log"))
	reflection.Register(s)

	srv, err := server.New(os.Getenv("recorder.addr"))
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}
	logger.Info().Msgf("Started at:%s", srv.Addr())
	err = srv.ServeGRPC(ctx, s)
	logger.Info().Msg("Stopped")
	return err
}
