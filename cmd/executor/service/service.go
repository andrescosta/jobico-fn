package service

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/process"
	"github.com/andrescosta/jobico/internal/executor"
)

const name = "executor"

type Setter func(*Service)

type Service struct {
	process.Container
	vm     *executor.VM
	dialer service.GrpcDialer
	option executor.Option
}

func New(ctx context.Context, ops ...Setter) (*Service, error) {
	s := &Service{
		option: executor.Option{},
		Container: process.Container{
			Name: name,
		},
	}
	for _, op := range ops {
		op(s)
	}
	vm, err := executor.NewVM(ctx, s.dialer, s.option)
	if err != nil {
		return nil, err
	}
	s.vm = vm
	svc, err := process.New(
		process.WithSidecarListener(s.ListenerOrDefault()),
		process.WithContext(ctx),
		process.WithName(name),
		process.WithAddr(s.AddrOrPanic()),
		process.WithHealthCheckFN(func(ctx context.Context) (map[string]string, error) {
			status := make(map[string]string)
			if !vm.IsUp() {
				return status, errors.New("error in executor")
			}
			return status, nil
		}),
		process.WithStarter(func(ctx context.Context) error {
			return vm.StartExecutors(ctx)
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

func (s *Service) Dispose() error {
	return s.vm.Close(s.Svc.Base.Ctx)
}

func WithGrpcDialer(d service.GrpcDialer) Setter {
	return func(s *Service) {
		s.dialer = d
	}
}

func WithHTTPConn(h service.HTTPConn) Setter {
	return func(s *Service) {
		s.Container.HTTPConn = h
	}
}

func WithOption(o executor.Option) Setter {
	return func(s *Service) {
		s.option = o
	}
}
