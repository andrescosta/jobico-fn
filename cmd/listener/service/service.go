package service

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/http"
	"github.com/andrescosta/jobico/internal/listener"
)

const name = "listener"

type Setter func(*Service)

type Service struct {
	http.Container
	controller    listener.Controller
	dialer        service.GrpcDialer
	listenerCache service.HTTPListener
}

func New(ctx context.Context, ops ...Setter) (*Service, error) {
	s := &Service{
		dialer:        service.DefaultGrpcDialer,
		listenerCache: service.DefaultGrpcListener,
		Container: http.Container{
			Name: name,
		},
	}
	for _, op := range ops {
		op(s)
	}
	_, _, err := env.Load(s.Name)
	if err != nil {
		return nil, err
	}
	c, err := listener.NewController(ctx, s.dialer, s.listenerCache)
	if err != nil {
		return nil, err
	}
	s.controller = c
	svc, err := http.New(
		http.WithListener[*http.ServiceOptions](s.ListenerOrDefault()),
		http.WithAddr(s.AddrOrPanic()),
		http.WithContext(ctx),
		http.WithName(name),
		http.WithHealthCheck[*http.ServiceOptions](func(ctx context.Context) (map[string]string, error) {
			return make(map[string]string), nil
		}),
		http.WithInitRoutesFn(c.ConfigureRoutes),
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
	return s.controller.Close()
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

func WithHTTPListener(l service.HTTPListener) Setter {
	return func(s *Service) {
		s.listenerCache = l
	}
}
