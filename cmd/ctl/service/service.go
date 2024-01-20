package service

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/ctl/server"
)

type Service struct {
	Listener service.GrpcListener
	DBOption *database.Option
}

func (s Service) Start(ctx context.Context) error {
	l := s.Listener
	if l == nil {
		l = service.DefaultGrpcListener
	}
	o := s.DBOption
	if o == nil {
		o = &database.Option{}
	}
	svc, err := grpc.New(
		grpc.WithName("ctl"),
		grpc.WithListener(l),
		grpc.WithContext(ctx),
		grpc.WithServiceDesc(&pb.Control_ServiceDesc),
		grpc.WithNewServiceFn(func(ctx context.Context) (any, error) {
			dbFileName := env.String("ctl.dbname", "db.db")
			return server.New(ctx, dbFileName, *o)
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
