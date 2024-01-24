package service

import (
	"context"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/queue"
	"github.com/andrescosta/jobico/internal/queue/controller"
)

type Service struct {
	Listener service.GrpcListener
	Dialer   service.GrpcDialer
	Option   *controller.Option
}

func (s Service) Start(ctx context.Context) error {
	l := s.Listener
	if l == nil {
		l = service.DefaultGrpcListener
	}
	d := s.Dialer
	if l == nil {
		l = service.DefaultGrpcDialer
	}
	o := s.Option
	if o == nil {
		o = &controller.Option{}
	}
	svc, err := grpc.New(
		grpc.WithListener(l),
		grpc.WithName("queue"),
		grpc.WithContext(ctx),
		grpc.WithServiceDesc(&pb.Queue_ServiceDesc),
		grpc.WithNewServiceFn(func(ctx context.Context) (any, error) {
			return queue.NewServer(ctx, d, *o)
		}),
	)
	if err != nil {
		return err
	}
	if err = svc.Serve(); err != nil {
		return err
	}
	return nil
}
