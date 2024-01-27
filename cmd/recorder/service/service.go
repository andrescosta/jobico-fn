package service

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/recorder/controller"
	"github.com/andrescosta/jobico/internal/recorder/server"
)

const name = "recorder"

type Service struct {
	Listener service.GrpcListener
	Option   *controller.Option
	Dialer   service.GrpcDialer
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
		grpc.WithName(name),
		grpc.WithContext(ctx),
		grpc.WithServiceDesc(&pb.Recorder_ServiceDesc),
		grpc.WithNewServiceFn(func(ctx context.Context) (any, error) {
			return server.New(ctx, ".\\log.log", *o)
		}),
	)
	if err != nil {
		return err
	}
	defer svc.Dispose()
	if err = svc.Serve(); err != nil {
		return err
	}
	return nil
}

func (s Service) Addr() *string {
	return env.StringOrNil(name + ".addr")
}

func (s Service) Kind() service.Kind {
	return grpc.Kind
}

func (s Service) CheckHealth(ctx context.Context) error {
	d := s.Dialer
	if d == nil {
		d = service.DefaultGrpcDialer
	}
	cli, err := grpc.NewHelthCheckClient(ctx, *s.Addr(), d)
	defer func() {
		_ = cli.Close()
	}()
	if err != nil {
		return err
	}
	return cli.CheckOk(ctx, name)
}
