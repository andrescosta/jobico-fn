package service

import (
	"context"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/http"
	"github.com/andrescosta/jobico/internal/listener"
)

type Service struct {
	Dialer   service.GrpcDialer
	Listener service.HTTPListener
}

func (s Service) Start(ctx context.Context) error {
	d := s.Dialer
	if d == nil {
		d = service.DefaultGrpcDialer
	}
	l := s.Listener
	if l == nil {
		l = service.DefaultHTTPListener
	}
	c, err := listener.New(ctx, d)
	if err != nil {
		return err
	}
	svc, err := http.New(
		http.WithListener[*http.ServiceOptions](l),
		http.WithContext(ctx),
		http.WithName("listener"),
		http.WithInitRoutesFn(c.ConfigureRoutes),
	)
	if err != nil {
		return err
	}
	if err := svc.Serve(); err != nil {
		return err
	}
	return nil
}
