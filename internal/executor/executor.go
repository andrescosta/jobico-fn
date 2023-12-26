package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/andrescosta/goico/pkg/broadcaster"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/wasm/wazero"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	cacheDir = "cache"
	NoError  = 0
)

type Executor struct {
	queue  string
	cancel context.CancelFunc
}
type VM struct {
	packages sync.Map
	w        *sync.WaitGroup
}
type jobPackage struct {
	PackageID string
	tenant    string
	Executors []*Executor
	Runtime   *wazero.WasmRuntime
	Modules   map[string]*wazero.WasmModuleString
	NextStep  map[string]*pb.ResultDef
}

func NewExecutorMachine(ctx context.Context) (*VM, error) {
	e := &VM{packages: sync.Map{}}
	if err := e.load(ctx); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *VM) StartExecutors(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	defer func() {
		e.packages.Range(func(key, value any) bool {
			pkg, _ := value.(*jobPackage)
			for _, m := range pkg.Modules {
				m.Close(ctx)
			}
			pkg.Runtime.Close(ctx)
			return true
		})
	}()
	var w sync.WaitGroup
	e.w = &w
	e.packages.Range(func(key, value any) bool {
		pkg := value.(*jobPackage)
		e.startPackage(ctx, pkg)
		return true
	})
	if err := e.startListeningUpdates(ctx); err != nil {
		logger.Warn().AnErr("error", err).Msg("updates are not being listened because an error")
	}
	logger.Info().Msg("Workers started")
	w.Wait()
	logger.Info().Msg("Workers stoped")
	return nil
}

func (e *VM) startPackage(ctx context.Context, pkg *jobPackage) {
	e.w.Add(1)
	for _, ex := range pkg.Executors {
		ctxE, cancel := context.WithCancel(ctx)
		ex.cancel = cancel
		go e.execute(ctxE, pkg.tenant, ex.queue, pkg.Modules, pkg.NextStep)
	}
}

func (e *VM) load(ctx context.Context) error {
	c, err := remote.NewControlClient(ctx)
	if err != nil {
		return err
	}
	ps, err := c.GetAllPackages(ctx)
	if err != nil {
		return err
	}
	if err := e.addOrUpdatePackages(ctx, ps); err != nil {
		return err
	}
	return nil
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
	modulesForEvents := make(map[string]*wazero.WasmModuleString)
	nextStepForEvents := make(map[string]*pb.ResultDef)
	r, err := wazero.NewWasmRuntime(ctx, env.InWorkDir(cacheDir))
	if err != nil {
		return err
	}
	jobPackage.Runtime = r
	executors := make([]*Executor, 0)
	for _, q := range pkg.Queues {
		executors = append(executors, &Executor{queue: q.ID})
	}
	jobPackage.Executors = executors
	repoClient, err := remote.NewRepoClient(ctx)
	if err != nil {
		return err
	}
	files := make(map[string][]byte)
	for _, job := range pkg.Jobs {
		event := job.Event
		for _, runtime := range pkg.Runtimes {
			if runtime.ID == event.Runtime {
				wasmfile, ok := files[runtime.ModuleRef]
				if !ok {
					wasmfile, err = repoClient.GetFile(ctx, pkg.Tenant, runtime.ModuleRef)
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
				wasmModule, err := wazero.NewWasmModuleString(ctx, runtime.ID, jobPackage.Runtime, wasmfile, funcName)
				if err != nil {
					return err
				}
				modulesForEvents[getModuleName(event.SupplierQueue, event.ID)] = wasmModule
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

func getFullPackageID(tenant string, queue string) string {
	return tenant + "/" + queue
}

func getModuleName(supplierQueue string, eventID string) string {
	return supplierQueue + "/" + eventID
}

func (e *VM) startListeningUpdates(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	controlClient, err := remote.NewControlClient(ctx)
	if err != nil {
		return err
	}
	l, err := controlClient.ListenerForPackageUpdates(ctx)
	if err != nil {
		return err
	}
	e.w.Add(1)
	go func() {
		err := e.listeningUpdates(ctx, l)
		if err != nil {
			logger.Info().AnErr("context.error", err).Msg("CTL update listener stopped.")
		}
	}()
	return nil
}

func (e *VM) listeningUpdates(ctx context.Context, l *broadcaster.Listener[*pb.UpdateToPackagesStrReply]) error {
	defer e.w.Done()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case u := <-l.C:
			e.onUpdate(ctx, u)
		}
	}
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
	o, _ := e.packages.Load(getFullPackageID(pkg.Tenant, pkg.ID))
	p := o.(*jobPackage)
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

func (e *VM) execute(ctx context.Context, tenant string, queue string, modules map[string]*wazero.WasmModuleString, nextSteps map[string]*pb.ResultDef) {
	defer e.w.Done()
	logger := zerolog.Ctx(ctx)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	logger.Debug().Msgf("Worker for Tenant: %s and queue: %s started", tenant, queue)
	queueErrors := 0
	maxQueueErrors := env.AsInt("max.queue.errors", 10)
	for {
		select {
		case <-ctx.Done():
			logger.Debug().Msgf("Worker for Tenant: %s and queue: %s stopped", tenant, queue)
			return
		case <-ticker.C:
			items, err := dequeue(ctx, tenant, queue)
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
				module, ok := modules[getModuleName(queue, item.Event)]
				if !ok {
					logger.Warn().Msgf("event %s not supported", item.Event)
					continue
				}
				code, result, err := executeWasm(ctx, module, item.Data)
				if err != nil {
					logger.Err(err).Msg("error executing")
				}
				if nextSteps, ok := nextSteps[item.Event]; ok {
					if err := reportToRecorder(ctx, queue, item.Event, tenant, code, result); err != nil {
						logger.Err(err).Msg("error reporting to recorder")
					}
					if err := makeDecisions(ctx, item.Event, tenant, code, result, nextSteps); err != nil {
						logger.Err(err).Msg("error enqueuing the result")
					}
				}
			}
		}
	}
}

func makeDecisions(ctx context.Context, _ string, tenant string, code uint64, _ string, resultDef *pb.ResultDef) error {
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
	client, err := remote.NewQueueClient(ctx)
	if err != nil {
		return err
	}
	if err := client.Queue(ctx, q); err != nil {
		return err
	}
	return nil
}

func reportToRecorder(ctx context.Context, queue string, eventID string, tenant string, code uint64, result string) error {
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
			Code:    code,
			Message: result,
		},
	}
	client, err := remote.NewRecorderClient(ctx)
	if err != nil {
		return err
	}
	return client.AddJobExecution(ctx, ex)
}

func dequeue(ctx context.Context, tenant string, queue string) ([]*pb.QueueItem, error) {
	client, err := remote.NewQueueClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Dequeue(ctx, tenant, queue)
}

func executeWasm(ctx context.Context, module *wazero.WasmModuleString, data []byte) (uint64, string, error) {
	mod := "goenv"
	logger := zerolog.Ctx(ctx)
	code, result, err := module.ExecuteMainFunc(ctx, string(data))
	if err != nil {
		return 0, "", errors.Join(err, fmt.Errorf("error in module %s", mod))
	}
	logger.Debug().Msgf("%d | %s", code, result)
	return code, result, nil
}
