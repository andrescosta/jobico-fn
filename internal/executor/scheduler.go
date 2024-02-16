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

type status int

const (
	statusStopped status = iota + 1
	statusStarting
	statusStarted
)

func newScheduller(ctx context.Context, ticker syncutil.Ticker) *scheduler {
	return &scheduler{
		currStatus: statusStarting,
		muStatus:   &sync.RWMutex{},
		ctx:        ctx,
		ticker:     ticker,
		executors:  collection.NewSyncMap[string, *process](),
	}
}

type scheduler struct {
	currStatus status
	muStatus   *sync.RWMutex
	ctx        context.Context
	ticker     syncutil.Ticker
	executors  *collection.SyncMap[string, *process]
}

func (s *scheduler) add(ex *process) {
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
	defer s.setCurrStatus(statusStopped)
	s.setCurrStatus(statusStarted)
	for {
		select {
		case <-s.ctx.Done():
			return
		case _, ok := <-s.ticker.Chan():
			if ok {
				w := sync.WaitGroup{}
				running := 0
				s.executors.Range(func(_ string, ex *process) bool {
					if s.ctx.Err() != nil {
						return false
					}
					w.Add(1)
					running++
					go ex.processEvents(s.ctx, &w)
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
	s.executors.Range(func(_ string, ex *process) bool {
		for _, e := range ex.events {
			err = errors.Join(e.module.wasmModule.Close(ctx), err)
		}
		return true
	})
	return err
}

func (s *scheduler) status() status {
	s.muStatus.RLock()
	defer s.muStatus.RUnlock()
	return s.currStatus
}

func (s *scheduler) setCurrStatus(status status) {
	s.muStatus.Lock()
	defer s.muStatus.Unlock()
	s.currStatus = status
}
