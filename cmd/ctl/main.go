package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/ctl/server"
)

func main() {
	svc, err := grpc.New(
		grpc.WithName("ctl"),
		grpc.WithContext(context.Background()),
		grpc.WithServiceDesc(&pb.Control_ServiceDesc),
		grpc.WithInitHandler(func(ctx context.Context) (any, error) {
			dbFileName := env.String("ctl.dbname", "db.db")
			return server.New(ctx, dbFileName)
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
