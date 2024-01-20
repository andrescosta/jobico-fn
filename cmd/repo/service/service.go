package service

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/repo/controller"
	"github.com/andrescosta/jobico/internal/repo/server"
)

type Service struct {
	Listener service.GrpcListener
	Option   *controller.Option
}

func (s Service) Start(ctx context.Context) error {
	l := s.Listener
	if l == nil {
		l = service.DefaultGrpcListener
	}
	o := s.Option
	if o == nil {
		o = &controller.Option{}
	}
	svc, err := grpc.New(
		grpc.WithListener(l),
		grpc.WithName("repo"),
		grpc.WithContext(ctx),
		grpc.WithServiceDesc(&pb.Repo_ServiceDesc),
		grpc.WithNewServiceFn(func(ctx context.Context) (any, error) {
			return server.New(ctx, env.String("repo.dir", "./"), *o), nil
		}),
	)
	if err != nil {
		log.Panicf("error starting repo service: %s", err)
	}

	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving repo service: %s", err)
	}
	return nil
}
