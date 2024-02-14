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

type Setter func(*Service)

type Service struct {
	grpc.Container
	option controller.Option
}

func New(ctx context.Context, ops ...Setter) (*Service, error) {
	s := &Service{
		option: controller.Option{},
		Container: grpc.Container{
			Name: name,
			GrpcConn: service.GrpcConn{
				Dialer:   service.DefaultGrpcDialer,
				Listener: service.DefaultGrpcListener,
			},
		},
	}
	for _, op := range ops {
		op(s)
	}
	_, _, err := env.Load(s.Name)
	if err != nil {
		return nil, err
	}
	svc, err := grpc.New(
		grpc.WithListener(s.Listener),
		grpc.WithAddr(s.AddrOrPanic()),
		grpc.WithName(s.Name),
		grpc.WithContext(ctx),
		grpc.WithServiceDesc(&pb.Recorder_ServiceDesc),
		grpc.WithProfilingEnabled(env.Bool("prof.enabled", false)),
		grpc.WithPProfAddr(env.StringOrNil("pprof.addr")),
		grpc.WithHealthCheckFn(func(_ context.Context) error { return nil }),
		grpc.WithNewServiceFn(func(ctx context.Context) (any, error) {
			return server.New(ctx, ".\\log.log", s.option)
		}),
	)
	if err != nil {
		return nil, err
	}
	s.Svc = svc
	return s, nil
}

func (s *Service) Start() error {
	return s.Svc.Serve()
}

func (s *Service) Dispose() {
	s.Svc.Dispose()
}

func WithOption(o controller.Option) Setter {
	return func(s *Service) {
		s.option = o
	}
}

func WithGrpcConn(g service.GrpcConn) Setter {
	return func(s *Service) {
		s.Container.GrpcConn = g
	}
}
