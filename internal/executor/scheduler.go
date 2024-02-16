package executor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/andrescosta/goico/pkg/collection"
	"github.com/andrescosta/goico/pkg/syncutil"
)

var maxExecutors = 10

func NewScheduller(ctx context.Context, ticker syncutil.Ticker) *scheduler {
	return &scheduler{
		statusScheduller: StatusStarting,
		muStatus:         &sync.RWMutex{},
		ctx:              ctx,
		ticker:           ticker,
		executors:        collection.NewSyncMap[string, *executor](),
	}
}

type scheduler struct {
	statusScheduller status
	muStatus         *sync.RWMutex
	ctx              context.Context
	ticker           syncutil.Ticker
	executors        *collection.SyncMap[string, *executor]
}

func (s *scheduler) add(ex *executor) {
	s.executors.Store(id(ex.tenant, ex.packageID, ex.queue), ex)
}

func id(tenant string, pkg string, queue string) string {
	return fmt.Sprintf("%s/%s/%s", tenant, pkg, queue)
}

func (s *scheduler) remove(tenant string, pkg string, queue string) {
	s.executors.Delete(id(tenant, pkg, queue))
}

func (s *scheduler) run() {
	defer s.ticker.Stop()
	defer s.setStatus(StatusStopped)
	s.setStatus(StatusStarted)
	for {
		select {
		case <-s.ctx.Done():
			return
		case _, ok := <-s.ticker.Chan():
			if ok {
				w := sync.WaitGroup{}
				running := 0
				s.executors.Range(func(id string, ex *executor) bool {
					if s.ctx.Err() != nil {
						return false
					}
					w.Add(1)
					running += 1
					go ex.execute(s.ctx, &w)
					if running == maxExecutors {
						w.Wait()
						running = 0
					}
					return true
				})
				if running > 0 {
					w.Wait()
				}
			}
		}
	}
}

func (s *scheduler) dispose() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.executors.Range(func(id string, ex *executor) bool {
		for _, e := range ex.events {
			err = errors.Join(e.module.wasmModule.Close(ctx), err)
		}
		return true
	})
	return err
}

func (e *scheduler) status() status {
	e.muStatus.RLock()
	defer e.muStatus.RUnlock()
	return e.statusScheduller
}

func (ex *scheduler) setStatus(status status) {
	ex.muStatus.Lock()
	defer ex.muStatus.Unlock()
	ex.statusScheduller = status
}
