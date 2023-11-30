package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/workflew/api/types"
	server1 "github.com/andrescosta/workflew/ctl/internal/server"
)

func main() {
	svc, err := service.NewGrpService(context.Background(), "ctl", &pb.Control_ServiceDesc,
		func(ctx context.Context) (any, error) {
			return server1.NewCotrolServer()
		})
	defer svc.Dispose()
	if err != nil {
		log.Panicf("error starting ctl service: %s", err)
	}
	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving ctl service %s", err)
	}
}
