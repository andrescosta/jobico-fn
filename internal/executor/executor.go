package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/execs/wasm"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/jobico/internal/api/client"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	cacheDir = "cache"
	NoError  = 0
)

type statusExecutor int

const (
	Stopped statusExecutor = iota + 1
	Starting
	Started
)

type Executor struct {
	mux        sync.RWMutex
	status     statusExecutor
	queue      string
	jobPackage *jobPackage
	tick       ticker
	cancel     context.CancelFunc
	vm         *VM
}
type VM struct {
	packages     sync.Map
	w            *sync.WaitGroup
	recorder     *client.Recorder
	ctl          *client.Ctl
	queue        *client.Queue
	repo         *client.Repo
	manualWakeup bool
}
type module struct {
	id         uint32
	wasmModule *wasm.Module
	tenant     string
	event      string
	vm         *VM
}
type jobPackage struct {
	PackageID string
	tenant    string
	Executors []*Executor
	Runtime   *wasm.Runtime
	Modules   map[string]module
	NextStep  map[string]*pb.ResultDef
}

type Option struct {
	ManualWakeup bool
}

func NewVM(ctx context.Context, d service.GrpcDialer, o Option) (*VM, error) {
	recorder, err := client.NewRecorder(ctx, d)
	if err != nil {
		return nil, err
	}
	ctl, err := client.NewCtl(ctx, d)
	if err != nil {
		return nil, err
	}
	queue, err := client.NewQueue(ctx, d)
	if err != nil {
		return nil, err
	}
	repo, err := client.NewRepo(ctx, d)
	if err != nil {
		return nil, err
	}
	e := &VM{
		packages:     sync.Map{},
		recorder:     recorder,
		manualWakeup: o.ManualWakeup,
		w:            &sync.WaitGroup{},
		ctl:          ctl,
		queue:        queue,
		repo:         repo,
	}

	return e, nil
}

func (e *VM) Close(ctx context.Context) error {
	var err error
	e.packages.Range(func(key, value any) bool {
		pkg, _ := value.(*jobPackage)
		for _, m := range pkg.Modules {
			err = errors.Join(err, m.wasmModule.Close(ctx))
		}
		err = errors.Join(err, pkg.Runtime.Close(ctx))
		return true
	})
	err = errors.Join(err, e.recorder.Close())
	err = errors.Join(err, e.ctl.Close())
	err = errors.Join(err, e.queue.Close())
	err = errors.Join(err, e.repo.Close())
	return err
}

func (e *VM) StartExecutors(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	if err := e.populate(ctx); err != nil {
		return err
	}

	e.packages.Range(func(key, value any) bool {
		pkg := value.(*jobPackage)
		e.startPackage(ctx, pkg)
		return true
	})
	logger.Info().Msg("Workers started")
	e.w.Wait()
	logger.Info().Msg("Workers stopped")
	return nil
}

func (e *VM) IsUp() bool {
	ok := true
	e.packages.Range(func(_, value any) bool {
		v := value.(*jobPackage)
		for _, p := range v.Executors {
			if p.getStatus() == Stopped {
				ok = false
				return false
			}
		}
		return true
	})
	return ok
}

func (ex *Executor) getStatus() statusExecutor {
	ex.mux.RLock()
	defer ex.mux.RUnlock()
	return ex.status
}

func (ex *Executor) setStatus(status statusExecutor) {
	ex.mux.Lock()
	defer ex.mux.Unlock()
	ex.status = status
}

func (e *VM) startPackage(ctx context.Context, pkg *jobPackage) {
	e.w.Add(len(pkg.Executors))
	for _, ex := range pkg.Executors {
		ctxE, cancel := context.WithCancel(ctx)
		ex.cancel = cancel
		ex.setStatus(Starting)
		go ex.execute(ctxE, e.w)
	}
}

func (e *VM) populate(ctx context.Context) error {
	ps, err := e.ctl.AllPackages(ctx)
	if err != nil {
		return err
	}
	if err := e.addOrUpdatePackages(ctx, ps); err != nil {
		return err
	}
	err = e.startListeningUpdates(ctx)
	return err
}

func (e *VM) addOrUpdatePackages(ctx context.Context, pkgs []*pb.JobPackage) error {
	for _, pkg := range pkgs {
		err := e.addPackage(ctx, pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *VM) addPackage(ctx context.Context, pkg *pb.JobPackage) error {
	jobPackage := &jobPackage{}
	jobPackage.tenant = pkg.Tenant
	jobPackage.PackageID = pkg.ID
	modulesForEvents := make(map[string]module)
	nextStepForEvents := make(map[string]*pb.ResultDef)
	r, err := wasm.NewRuntime(env.WorkdirPlus(cacheDir))
	if err != nil {
		return err
	}
	jobPackage.Runtime = r
	executors := make([]*Executor, 0)

	for _, q := range pkg.Queues {
		if strings.HasSuffix(q.ID, "_ok") || strings.HasSuffix(q.ID, "_error") {
			continue
		}
		var tick ticker
		if e.manualWakeup {
			tick = &channelBasedTicker{
				c: make(chan time.Time),
			}
		} else {
			dur := *env.Duration("executor.timeout", 5*time.Second)
			tick = &timeBasedTicker{
				ticker: time.NewTicker(dur),
			}
		}
		executors = append(executors,
			&Executor{
				queue:      q.ID,
				jobPackage: jobPackage,
				vm:         e,
				mux:        sync.RWMutex{},
				status:     Stopped,
				tick:       tick,
			})
	}
	jobPackage.Executors = executors
	files := make(map[string][]byte)
	var id uint32
	for _, job := range pkg.Jobs {
		event := job.Event
		for _, runtime := range pkg.Runtimes {
			if runtime.ID == event.Runtime {
				wasmfile, ok := files[runtime.ModuleRef]
				if !ok {
					wasmfile, err = e.repo.File(ctx, pkg.Tenant, runtime.ModuleRef)
					if err != nil {
						return err
					}
					files[runtime.ModuleRef] = wasmfile
				}
				// We create a module per queue because wazero module call is not goroutine compatible
				funcName := "event"
				if runtime.MainFuncName != nil {
					funcName = *runtime.MainFuncName
				}
				module := module{
					id:     id,
					event:  event.ID,
					tenant: pkg.Tenant,
					vm:     e,
				}
				id = id + 1
				wasmModule, err := wasm.NewModule(ctx, jobPackage.Runtime, wasmfile, funcName, module.sendLogToRecorder)
				if err != nil {
					return err
				}
				module.wasmModule = wasmModule
				modulesForEvents[getModuleName(event.SupplierQueue, event.ID)] = module
				break
			}
		}
		if job.Result != nil {
			nextStepForEvents[event.ID] = job.Result
		}
	}
	jobPackage.Modules = modulesForEvents
	jobPackage.NextStep = nextStepForEvents
	e.packages.Store(getFullPackageID(pkg.Tenant, pkg.ID), jobPackage)
	return nil
}

func (m module) sendLogToRecorder(ctx context.Context, _ uint32, lvl uint32, msg string) error {
	now := time.Now()
	host, err := os.Hostname()
	if err != nil {
		host = "<error>"
	}

	return m.vm.recorder.AddJobExecution(ctx, &pb.JobExecution{
		Event:  m.event,
		Tenant: m.tenant,
		Queue:  "",
		Date: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		Server: host,
		Result: &pb.JobResult{
			Type:     pb.JobResult_Log,
			TypeDesc: "log",
			Code:     uint64(lvl),
			Message:  msg,
		},
	})
}

func getFullPackageID(tenant string, queue string) string {
	return tenant + "/" + queue
}

func getModuleName(supplierQueue string, eventID string) string {
	return supplierQueue + "/" + eventID
}

func (e *VM) startListeningUpdates(ctx context.Context) error {
	l, err := e.ctl.ListenerForPackageUpdates(ctx)
	if err != nil {
		return err
	}
	e.w.Add(1)
	go func() {
		defer e.w.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case u := <-l.C:
				e.onUpdate(ctx, u)
			}
		}
	}()
	return nil
}

func (e *VM) onUpdate(ctx context.Context, u *pb.UpdateToPackagesStrReply) {
	logger := zerolog.Ctx(ctx)
	switch u.Type {
	case pb.UpdateType_New:
		if err := e.newPackage(ctx, u.Object); err != nil {
			logger.Warn().AnErr("error", err).Msg("onUpdate: error adding package")
		}
	case pb.UpdateType_Update:
		if err := e.updatePackage(ctx, u.Object); err != nil {
			logger.Warn().AnErr("error", err).Msg("onUpdate: error adding package")
		}
	case pb.UpdateType_Delete:
		e.deletePackage(ctx, u.Object)
	}
}

func (e *VM) newPackage(ctx context.Context, pkg *pb.JobPackage) error {
	if err := e.addPackage(ctx, pkg); err != nil {
		return err
	}
	j, _ := e.packages.Load(getFullPackageID(pkg.Tenant, pkg.ID))
	p := j.(*jobPackage)
	e.startPackage(ctx, p)
	return nil
}

func (e *VM) deletePackage(_ context.Context, p *pb.JobPackage) {
	o, ok := e.packages.LoadAndDelete(getFullPackageID(p.Tenant, p.ID))
	if ok {
		pkg := o.(*jobPackage)
		for _, e := range pkg.Executors {
			e.cancel()
		}
	}
}

func (e *VM) updatePackage(ctx context.Context, p *pb.JobPackage) error {
	e.deletePackage(ctx, p)
	return e.newPackage(ctx, p)
}

func (ex *Executor) execute(ctx context.Context, w *sync.WaitGroup) {
	defer w.Done()
	defer ex.setStatus(Stopped)
	defer ex.tick.Stop()
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msgf("Worker for Tenant: %s and queue: %s started", ex.jobPackage.tenant, ex.queue)
	queueErrors := 0
	maxQueueErrors := env.Int("max.queue.errors", 10)
	chTick := ex.tick.Chan()
	ex.setStatus(Started)
	for {
		select {
		case <-ctx.Done():
			logger.Debug().Msgf("Worker for Tenant: %s and queue: %s stopped", ex.jobPackage.tenant, ex.queue)
			return
		case <-chTick:
			items, err := ex.vm.queue.Dequeue(ctx, ex.jobPackage.tenant, ex.queue)
			if err != nil {
				queueErrors++
				logger.Err(err).Msg("error dequeuing")
				if queueErrors == maxQueueErrors {
					logger.Err(err).Msg("Queue: stopping goroutine due to many errors")
					return
				}
				continue
			}
			queueErrors = 0
			for _, item := range items {
				module, ok := ex.jobPackage.Modules[getModuleName(ex.queue, item.Event)]
				if !ok {
					logger.Warn().Msgf("event %s not supported", item.Event)
					continue
				}
				code, result, err := executeWasm(ctx, module.wasmModule, module.id, item.Data)
				if err != nil {
					logger.Err(err).Msg("error executing")
				}
				if err := ex.vm.reportResultToRecorder(ctx, ex.queue, item.Event, ex.jobPackage.tenant, code, result); err != nil {
					logger.Err(err).Msg("error reporting to recorder")
				}
				if nextSteps, ok := ex.jobPackage.NextStep[item.Event]; ok {
					if err := ex.vm.makeDecisions(ctx, item.Event, ex.jobPackage.tenant, code, result, nextSteps); err != nil {
						logger.Err(err).Msg("error enqueuing the result")
					}
				}
			}
		}
	}
}

func (e *VM) makeDecisions(ctx context.Context, _ string, tenant string, code uint64, _ string, resultDef *pb.ResultDef) error {
	r := pb.JobResult{
		Code: code,
	}
	bytes1, err := proto.Marshal(&r)
	if err != nil {
		return err
	}
	var q *pb.QueueRequest
	if code == NoError {
		q = &pb.QueueRequest{
			Tenant: tenant,
			Queue:  resultDef.Ok.SupplierQueue,
			Items: []*pb.QueueItem{
				{
					Event: resultDef.Ok.ID,
					Data:  bytes1,
				},
			},
		}
	} else {
		q = &pb.QueueRequest{
			Tenant: tenant,
			Queue:  resultDef.Error.SupplierQueue,
			Items: []*pb.QueueItem{
				{
					Event: resultDef.Error.ID,
					Data:  bytes1,
				},
			},
		}
	}
	if err := e.queue.Queue(ctx, q); err != nil {
		return err
	}
	return nil
}

func (e *VM) reportResultToRecorder(ctx context.Context, queue string, eventID string, tenant string, code uint64, result string) error {
	now := time.Now()
	host, err := os.Hostname()
	if err != nil {
		host = "<error>"
	}
	ex := &pb.JobExecution{
		Event:  eventID,
		Tenant: tenant,
		Queue:  queue,
		Date: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		Server: host,
		Result: &pb.JobResult{
			Type:     pb.JobResult_Result,
			TypeDesc: "result",
			Code:     code,
			Message:  result,
		},
	}
	return e.recorder.AddJobExecution(ctx, ex)
}

func executeWasm(ctx context.Context, module *wasm.Module, id uint32, data []byte) (uint64, string, error) {
	mod := "goenv"
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, *env.Duration("wasm.exec.timeout", 2*time.Minute))
	defer cancel()
	code, result, err := module.Execute(ctx, id, string(data))
	if err != nil {
		return 0, "", errors.Join(err, fmt.Errorf("error in module %s", mod))
	}
	logger.Debug().Msgf("%d | %s", code, result)
	return code, result, nil
}
