package service

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/ctl/server"
)

const name = "ctl"

type Setter func(*Service)

type Service struct {
	grpc.Container
	dbOption database.Option
}

func New(ctx context.Context, ops ...Setter) (*Service, error) {
	s := &Service{
		dbOption: database.Option{},
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
		grpc.WithName(s.Name),
		grpc.WithAddr(s.AddrOrPanic()),
		grpc.WithListener(s.Listener),
		grpc.WithContext(ctx),
		grpc.WithHealthCheckFn(func(_ context.Context) error { return nil }),
		grpc.WithServiceDesc(&pb.Control_ServiceDesc),
		grpc.WithProfilingEnabled(env.Bool("prof.enabled", false)),
		grpc.WithPProfAddr(env.StringOrNil("pprof.addr")),
		grpc.WithNewServiceFn(func(ctx context.Context) (any, error) {
			dbDir := env.String("ctl.dbdir", "db")
			return server.New(ctx, dbDir, s.dbOption)
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

func WithDBOption(d database.Option) Setter {
	return func(s *Service) {
		s.dbOption = d
	}
}

func WithGrpcConn(g service.GrpcConn) Setter {
	return func(s *Service) {
		s.Container.GrpcConn = g
	}
}
