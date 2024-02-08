package test

import (
	"context"
	"errors"
	"time"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	ctl "github.com/andrescosta/jobico/cmd/ctl/service"
	exec "github.com/andrescosta/jobico/cmd/executor/service"
	listener "github.com/andrescosta/jobico/cmd/listener/service"
	queue "github.com/andrescosta/jobico/cmd/queue/service"
	recorder "github.com/andrescosta/jobico/cmd/recorder/service"
	repo "github.com/andrescosta/jobico/cmd/repo/service"
	"github.com/andrescosta/jobico/internal/executor"
	queuectl "github.com/andrescosta/jobico/internal/queue/controller"
	recorderctl "github.com/andrescosta/jobico/internal/recorder/controller"
	repoctl "github.com/andrescosta/jobico/internal/repo/controller"
)

type jobicoPlatform struct {
	conn     *service.BufConn
	ctl      *ctl.Service
	queue    *queue.Service
	repo     *repo.Service
	listener *listener.Service
	executor *exec.Service
	recorder *recorder.Service
}

func (j *jobicoPlatform) dispose() error {
	errs := make([]error, 0)
	j.ctl.Dispose()
	j.queue.Dispose()
	err := j.listener.Dispose()
	if err != nil {
		errs = append(errs, err)
	}
	err = j.executor.Dispose()
	if err != nil {
		errs = append(errs, err)
	}
	j.recorder.Dispose()
	return errors.Join(errs...)
}

func newPlatform(ctx context.Context) (*jobicoPlatform, error) {
	return newPlatformWithTimeout(ctx, *env.Duration("dial.timeout"))
}

func newPlatformWithTimeout(ctx context.Context, time time.Duration) (*jobicoPlatform, error) {
	conn := service.NewBufConnWithTimeout(time)
	ctl, err := ctl.New(ctx,
		ctl.WithGrpcConn(service.GrpcConn{
			Listener: conn,
			Dialer:   conn,
		}),
		ctl.WithDBOption(database.Option{InMemory: true}))
	if err != nil {
		return nil, err
	}
	queue, err := queue.New(ctx, queue.WithGrpcConn(
		service.GrpcConn{
			Listener: conn,
			Dialer:   conn,
		}), queue.WithOption(queuectl.Option{InMemory: true}))
	if err != nil {
		return nil, err
	}
	repo, err := repo.New(ctx, repo.WithGrpcConn(
		service.GrpcConn{
			Listener: conn,
			Dialer:   conn,
		}), repo.WithOption(repoctl.Option{InMemory: true}))
	if err != nil {
		return nil, err
	}

	listener, err := listener.New(ctx, listener.WithHTTPConn(service.HTTPConn{
		ClientBuilder: conn,
		Listener:      conn,
	}), listener.WithGrpcDialer(conn), listener.WithHTTPListener(conn))
	if err != nil {
		return nil, err
	}

	executor, err := exec.New(ctx, exec.WithHTTPConn(service.HTTPConn{
		ClientBuilder: conn,
		Listener:      conn,
	}),
		exec.WithGrpcDialer(conn),
		exec.WithOption(
			executor.Option{ManualWakeup: false}))
	if err != nil {
		return nil, err
	}

	recorder, err := recorder.New(ctx,
		recorder.WithGrpcConn(service.GrpcConn{
			Listener: conn,
			Dialer:   conn,
		}), recorder.WithOption(recorderctl.Option{InMemory: true}))
	if err != nil {
		return nil, err
	}
	return &jobicoPlatform{
		ctl:      ctl,
		conn:     conn,
		queue:    queue,
		repo:     repo,
		listener: listener,
		executor: executor,
		recorder: recorder,
	}, nil
}
