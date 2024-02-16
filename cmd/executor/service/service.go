package service

import (
	"context"
	"errors"
	"time"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/process"
	"github.com/andrescosta/jobico/internal/executor"
)

const name = "executor"

type Setter func(*Service)

type Service struct {
	process.Container
	delay  time.Duration
	vm     *executor.VM
	dialer service.GrpcDialer
	option executor.Options
}

func New(ctx context.Context, ops ...Setter) (*Service, error) {
	s := &Service{
		option: executor.Options{},
		dialer: service.DefaultGrpcDialer,
		Container: process.Container{
			Name: name,
		},
		delay: 0,
	}
	for _, op := range ops {
		op(s)
	}
	_, _, err := env.Load(s.Name)
	if err != nil {
		return nil, err
	}
	s.delay = *env.Duration("executor.delay", 0)
	vm, err := executor.NewVM(ctx, s.dialer, s.option)
	if err != nil {
		return nil, err
	}
	s.vm = vm
	empty := make(map[string]string)
	svc, err := process.New(
		process.WithSidecarListener(s.ListenerOrDefault()),
		process.WithContext(ctx),
		process.WithName(name),
		process.WithAddr(s.AddrOrPanic()),
		process.WithProfilingEnabled(env.Bool("prof.enabled", false)),
		process.WithHealthCheckFN(func(_ context.Context) (map[string]string, error) {
			if !vm.IsUp() {
				return empty, errors.New("error in executor")
			}
			return empty, nil
		}),
		process.WithStarter(func(ctx context.Context) error {
			return vm.Start(ctx)
		}),
	)
	if err != nil {
		return nil, err
	}
	s.Svc = svc
	return s, nil
}

func (s *Service) Start() error {
	return s.Svc.ServeWithDelay(s.delay)
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

func WithOption(o executor.Options) Setter {
	return func(s *Service) {
		s.option = o
	}
}
