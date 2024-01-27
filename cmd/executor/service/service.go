package service

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/http"
	"github.com/andrescosta/goico/pkg/service/process"
	"github.com/andrescosta/jobico/internal/executor"
)

const name = "executor"

type Service struct {
	Dialer        service.GrpcDialer
	Listener      service.HTTPListener
	ClientBuilder service.HTTPClient
	Option        *executor.Option
}

func (s Service) Start(ctx context.Context) (err error) {
	d := s.Dialer
	if d == nil {
		d = service.DefaultGrpcDialer
	}
	l := s.Listener
	if l == nil {
		l = service.DefaultHTTPListener
	}
	o := s.Option
	if o == nil {
		o = &executor.Option{}
	}
	m, err := executor.NewVM(ctx, d, *o)
	if err != nil {
		return
	}
	defer func() {
		errc := m.Close(ctx)
		err = errors.Join(err, errc)
	}()
	e, err := process.New(
		process.WithSidecarListener(l),
		process.WithContext(ctx),
		process.WithName(name),
		process.WithHealthCheckFN(func(ctx context.Context) (map[string]string, error) {
			status := make(map[string]string)
			if !m.IsUp() {
				return status, errors.New("error in executor")
			}
			return status, nil
		}),
		process.WithServeHandler(func(ctx context.Context) error {
			m.StartExecutors(ctx)
			return nil
		}),
	)
	if err != nil {
		return
	}
	err = e.Serve()
	return
}

func (s Service) Addr() *string {
	return env.StringOrNil(name + ".addr")
}

func (s Service) Kind() service.Kind {
	return process.Kind
}

func (s Service) CheckHealth(ctx context.Context) error {
	b := s.ClientBuilder
	if b == nil {
		b = service.DefaultHTTPClient
	}
	cli, err := b.NewHTTPClient(*s.Addr())
	if err != nil {
		return err
	}
	return http.CheckServiceHealth(ctx, cli, *s.Addr())
}
