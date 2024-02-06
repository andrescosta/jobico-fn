package service

import (
	"context"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/queue"
	"github.com/andrescosta/jobico/internal/queue/controller"
)

const name = "queue"

type Setter func(*Service)

type Service struct {
	grpc.Container
	option controller.Option
}

func New(ctx context.Context, ops ...Setter) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	s := &Service{
		option: controller.Option{},
		Container: grpc.Container{
			Name:   name,
			Cancel: cancel,
			GrpcConn: service.GrpcConn{
				Dialer:   service.DefaultGrpcDialer,
				Listener: service.DefaultGrpcListener,
			},
		},
	}
	for _, op := range ops {
		op(s)
	}
	svc, err := grpc.New(
		grpc.WithListener(s.Listener),
		grpc.WithName(s.Name),
		grpc.WithAddr(s.AddrOrPanic()),
		grpc.WithContext(ctx),
		grpc.WithServiceDesc(&pb.Queue_ServiceDesc),
		grpc.WithHealthCheckFn(func(ctx context.Context) error { return nil }),
		grpc.WithNewServiceFn(func(ctx context.Context) (any, error) {
			return queue.NewServer(ctx, s.Dialer, s.option)
		}),
	)
	if err != nil {
		return nil, err
	}
	s.Svc = svc
	return s, nil
}

func (s *Service) Start() error {
	defer s.dispose()
	return s.Svc.Serve()
}

func (s *Service) dispose() {
	s.Svc.Dispose()
}

func (s *Service) Stop() {
	s.Cancel()
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
