package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/server"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/workflew/api/types"
	server1 "github.com/andrescosta/workflew/ctl/internal/server"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	service.StartNamed("CTL", serviceFunc)
}

func serviceFunc(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	s := grpc.NewServer()
	svr, err := server1.NewCotrolServer(ctx)
	if err != nil {
		return err
	}
	defer svr.Close(ctx)
	pb.RegisterControlServer(s, svr)
	reflection.Register(s)

	srv, err := server.New(os.Getenv("ctl.port"))
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}
	logger.Info().Msgf("Server started at:%s", srv.Addr())
	err = srv.ServeGRPC(ctx, s)
	logger.Info().Msg("Server stopped")
	return err
}
