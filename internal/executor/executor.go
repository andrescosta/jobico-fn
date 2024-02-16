package executor

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/execs/wasm"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/syncutil"
	"github.com/andrescosta/jobico/internal/api/client"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/rs/zerolog"
)

const (
	cacheDir = "cache"
	NoError  = 0
)

type cli struct {
	recorder *client.Recorder
	ctl      *client.Ctl
	queue    *client.Queue
	repo     *client.Repo
}

type Executor struct {
	cli       *cli
	scheduler *scheduler
	runtime   *wasm.Runtime
}

type Options struct {
	Ticker syncutil.Ticker
}

func New(ctx context.Context, dialer service.GrpcDialer, option Options) (*Executor, error) {
	recorder, err := client.NewRecorder(ctx, dialer)
	if err != nil {
		return nil, err
	}
	ctl, err := client.NewCtl(ctx, dialer)
	if err != nil {
		return nil, err
	}
	queue, err := client.NewQueue(ctx, dialer)
	if err != nil {
		return nil, err
	}
	repo, err := client.NewRepo(ctx, dialer)
	if err != nil {
		return nil, err
	}
	ticker := option.Ticker
	if ticker == nil {
		dur := *env.Duration("executor.timeout", 5*time.Second)
		ticker = &syncutil.TimeTicker{
			Ticker: time.NewTicker(dur),
		}
	}
	cli := &cli{
		recorder: recorder,
		ctl:      ctl,
		queue:    queue,
		repo:     repo,
	}
	wasmRuntime, err := wasm.NewRuntimeWithCompilationCache(env.WorkdirPlus(cacheDir))
	if err != nil {
		return nil, err
	}
	scheduller := newScheduller(ctx, ticker)
	e := &Executor{
		cli:       cli,
		scheduler: scheduller,
		runtime:   wasmRuntime,
	}
	return e, nil
}

func (e *Executor) Close(ctx context.Context) error {
	var err error
	err = errors.Join(e.runtime.Close(ctx))
	err = errors.Join(err, e.scheduler.dispose())
	err = errors.Join(err, e.cli.recorder.Close())
	err = errors.Join(err, e.cli.ctl.Close())
	err = errors.Join(err, e.cli.queue.Close())
	err = errors.Join(err, e.cli.repo.Close())
	return err
}

func (e *Executor) Start(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	if err := e.init(ctx); err != nil {
		return err
	}
	logger.Info().Msg("Workers started")
	e.scheduler.run()
	logger.Info().Msg("Workers stopped")
	return nil
}

func (e *Executor) IsUp() bool {
	return e.scheduler.status() == statusStarted
}

func (e *Executor) init(ctx context.Context) error {
	ps, err := e.cli.ctl.AllPackages(ctx)
	if err != nil {
		return err
	}
	for _, pkg := range ps {
		err := e.addExecutors(ctx, pkg)
		if err != nil {
			return err
		}
	}
	err = e.startListeningUpdates(ctx)
	return err
}

func (e *Executor) addExecutors(ctx context.Context, pkg *pb.JobPackage) error {
	events := make(map[string]*event)
	var id uint32
	for _, job := range pkg.Jobs {
		runtime := getRuntime(job.Event.Runtime, pkg.Runtimes)
		if runtime != nil {
			sender := &recorder{
				cli:    e.cli,
				tenant: pkg.Tenant,
				event:  job.Event.ID,
			}
			event := &event{
				id:        job.Event.ID,
				nextStep:  job.Result,
				logSender: sender,
			}
			events[job.Event.ID] = event
			wasmfile, err := e.cli.repo.File(ctx, pkg.Tenant, runtime.ModuleRef)
			if err != nil {
				return err
			}
			funcName := "event"
			if runtime.MainFuncName != nil {
				funcName = *runtime.MainFuncName
			}
			wasmModule, err := wasm.NewModule(ctx, e.runtime, wasmfile, funcName, sender.sendLog)
			if err != nil {
				return err
			}
			module := module{
				id:         id,
				wasmModule: wasmModule,
			}
			id = id + 1
			event.module = &module
		}
		for _, q := range pkg.Queues {
			if strings.HasSuffix(q.ID, "_ok") || strings.HasSuffix(q.ID, "_error") {
				continue
			}
			ex := &process{
				packageID: pkg.GetID(),
				tenant:    pkg.Tenant,
				queue:     q.ID,
				runtime:   e.runtime,
				events:    events,
				cli:       e.cli,
			}
			e.scheduler.add(ex)
		}
	}
	return nil
}

func (e *Executor) removeExecutor(_ context.Context, p *pb.JobPackage) {
	for _, q := range p.Queues {
		e.scheduler.remove(p.Tenant, p.ID, q.ID)
	}
}

func (e *Executor) startListeningUpdates(ctx context.Context) error {
	l, err := e.cli.ctl.ListenerForPackageUpdates(ctx)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case u, ok := <-l.C:
				// todo: reconnect if not ok
				if ok {
					e.onUpdate(ctx, u)
				}
			}
		}
	}()
	return nil
}

func (e *Executor) onUpdate(ctx context.Context, u *pb.UpdateToPackagesStrReply) {
	logger := zerolog.Ctx(ctx)
	switch u.Type {
	case pb.UpdateType_New:
		if err := e.addExecutors(ctx, u.Object); err != nil {
			logger.Warn().AnErr("error", err).Msg("onUpdate: error adding package")
		}
	case pb.UpdateType_Update:
		e.removeExecutor(ctx, u.Object)
		if err := e.addExecutors(ctx, u.Object); err != nil {
			logger.Warn().AnErr("error", err).Msg("onUpdate: error adding package")
		}
	case pb.UpdateType_Delete:
		e.removeExecutor(ctx, u.Object)
	}
}

func getRuntime(id string, rs []*pb.RuntimeDef) *pb.RuntimeDef {
	for _, r := range rs {
		if r.ID == id {
			return r
		}
	}
	return nil
}
