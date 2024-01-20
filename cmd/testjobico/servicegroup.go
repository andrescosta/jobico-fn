package testjobico

import (
	"context"
	"errors"
	"sync"

	"github.com/andrescosta/goico/pkg/collection"
)

type Starter interface {
	Start(context.Context) error
}

type ServiceGroup struct {
	w      *sync.WaitGroup
	cancel context.CancelFunc
	ctx    context.Context
	qerrs  *collection.SyncQueue[error]
}

func NewServiceGroup() *ServiceGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &ServiceGroup{
		w:      &sync.WaitGroup{},
		cancel: cancel,
		ctx:    ctx,
		qerrs:  collection.NewQueue[error](),
	}
}

func (s *ServiceGroup) AddAndStart(services []Starter) {
	s.w.Add(len(services))
	for _, service := range services {
		go func(service Starter) {
			defer s.w.Done()
			if err := service.Start(s.ctx); err != nil {
				s.qerrs.Queue(err)
			}
		}(service)
	}
}

func (s *ServiceGroup) Stop() error {
	s.cancel()
	s.w.Wait()
	if s.qerrs.Size() > 0 {
		return errors.Join(s.qerrs.Slice()...)
	}
	return nil
}
