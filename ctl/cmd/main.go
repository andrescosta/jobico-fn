package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/ctl/internal/server"
)

func main() {
	svc, err := grpc.New(
		grpc.WithName("ctl"),
		grpc.WithContext(context.Background()),
		grpc.WithServiceDesc(&pb.Control_ServiceDesc),
		grpc.WithInitHandler(func(ctx context.Context) (any, error) {
			dbPath := env.Env("ctl.dbpath", ".\\db.db")
			return server.New(ctx, dbPath)
		}),
	)
	if err != nil {
		log.Panicf("error starting ctl service: %s", err)
	}
	defer svc.Dispose()
	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving ctl service %s", err)
	}
}
