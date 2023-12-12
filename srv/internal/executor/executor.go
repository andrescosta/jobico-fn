package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/wasmico/wazero"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	cacheDir = ".\\cache"
	NO_ERROR = 0
)

type Executor struct {
	queueId string
	cancel  context.CancelFunc
}

type ExecutorMachine struct {
	packages sync.Map
	w        *sync.WaitGroup
}

type jobPackage struct {
	PackageId string
	TenantId  string
	Executors []*Executor
	Runtime   *wazero.WasmRuntime
	Modules   map[string]*wazero.WasmModuleString
	NextStep  map[string]*pb.ResultDef
}

func NewExecutorMachine(ctx context.Context) (*ExecutorMachine, error) {
	e := &ExecutorMachine{packages: sync.Map{}}
	if err := e.load(ctx); err != nil {
		return nil, err
	}
	return e, nil

}

func (e *ExecutorMachine) StartExecutors(ctx context.Context) error {
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

	e.startListeningUpdates(ctx)
	logger.Info().Msg("Workers started")
	w.Wait()
	logger.Info().Msg("Workers stoped")
	return nil
}

func (e *ExecutorMachine) startPackage(ctx context.Context, pkg *jobPackage) {
	e.w.Add(1)
	for _, ex := range pkg.Executors {
		ctxE, cancel := context.WithCancel(ctx)
		ex.cancel = cancel
		go e.execute(ctxE, pkg.TenantId, ex.queueId, pkg.Modules, pkg.NextStep)
	}
}

func (e *ExecutorMachine) load(ctx context.Context) error {
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

func (e *ExecutorMachine) addOrUpdatePackages(ctx context.Context, pkgs []*pb.JobPackage) error {
	for _, pkg := range pkgs {
		err := e.addPackage(ctx, pkg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (e *ExecutorMachine) addPackage(ctx context.Context, pkg *pb.JobPackage) error {
	jobPackage := &jobPackage{}
	jobPackage.TenantId = pkg.TenantId
	jobPackage.PackageId = pkg.ID

	modulesForEvents := make(map[string]*wazero.WasmModuleString)
	nextStepForEvents := make(map[string]*pb.ResultDef)
	r, err := wazero.NewWasmRuntime(ctx, cacheDir)
	if err != nil {
		return err
	}
	jobPackage.Runtime = r
	executors := make([]*Executor, 0)
	for _, q := range pkg.Queues {
		executors = append(executors, &Executor{queueId: q.ID})
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
			if runtime.ID == event.RuntimeId {
				wasmfile, ok := files[runtime.ModuleRef]
				if !ok {
					wasmfile, err = repoClient.GetFile(ctx, pkg.TenantId, runtime.ModuleRef)
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
				modulesForEvents[getModuleName(event.SupplierQueueId, event.ID)] = wasmModule
				break
			}
		}
		if job.Result != nil {
			nextStepForEvents[event.ID] = job.Result
		}
	}
	jobPackage.Modules = modulesForEvents
	jobPackage.NextStep = nextStepForEvents
	e.packages.Store(getFullPackageId(pkg.TenantId, pkg.ID), jobPackage)
	return nil
}

func getFullPackageId(tenantId string, queueId string) string {
	return tenantId + "/" + queueId
}

func getModuleName(supplierQueueId string, eventId string) string {
	return supplierQueueId + "/" + eventId

}

func (j *ExecutorMachine) startListeningUpdates(ctx context.Context) error {
	controlClient, err := remote.NewControlClient(ctx)
	if err != nil {
		return err
	}
	l, err := controlClient.ListenerForPackageUpdates(ctx)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case u := <-l.C:
				j.onUpdate(ctx, u)
			}
		}
	}()
	return nil
}

func (j *ExecutorMachine) onUpdate(ctx context.Context, u *pb.UpdateToPackagesStrReply) {
	switch u.Type {
	case pb.UpdateType_New:
		j.newPackage(ctx, u.Object)
	case pb.UpdateType_Update:
		j.updatePackage(ctx, u.Object)
	case pb.UpdateType_Delete:
		j.deletePackage(ctx, u.Object)

	}
}
func (j *ExecutorMachine) newPackage(ctx context.Context, pkg *pb.JobPackage) {
	j.addPackage(ctx, pkg)
	o, _ := j.packages.Load(getFullPackageId(pkg.TenantId, pkg.ID))
	p := o.(*jobPackage)
	j.startPackage(ctx, p)
}

func (j *ExecutorMachine) deletePackage(ctx context.Context, p *pb.JobPackage) {
	o, ok := j.packages.LoadAndDelete(getFullPackageId(p.TenantId, p.ID))
	if ok {
		pkg := o.(*jobPackage)
		for _, e := range pkg.Executors {
			e.cancel()
		}
	}
}

func (j *ExecutorMachine) updatePackage(ctx context.Context, p *pb.JobPackage) {
	j.deletePackage(ctx, p)
	j.newPackage(ctx, p)
}

func (e *ExecutorMachine) execute(ctx context.Context, tenantId string, queueId string, modules map[string]*wazero.WasmModuleString, nextSteps map[string]*pb.ResultDef) {
	defer e.w.Done()
	logger := zerolog.Ctx(ctx)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	logger.Debug().Msgf("Worker for tenant: %s and queue: %s started", tenantId, queueId)
	queueErrors := 0
	maxQueueErrors := env.GetAsInt("max.queue.errors", 10)
	for {
		select {
		case <-ctx.Done():
			logger.Debug().Msgf("Worker for tenant: %s and queue: %s stopped", tenantId, queueId)
			return
		case <-ticker.C:
			items, err := dequeue(ctx, tenantId, queueId)
			if err != nil {
				queueErrors++
				logger.Err(err).Msg("error dequeuing")
				if queueErrors == maxQueueErrors {
					logger.Err(err).Msg("stopping goroutine due queue errors")
					return
				} else {
					continue
				}
			}
			queueErrors = 0
			for _, item := range items {
				module, ok := modules[getModuleName(queueId, item.EventId)]
				if !ok {
					logger.Warn().Msgf("event %s not supported", item.EventId)
					continue
				}
				code, result, err := executeWasm(ctx, module, item.Data)
				if err != nil {
					logger.Err(err).Msg("error executing")
				}
				if nextSteps, ok := nextSteps[item.EventId]; ok {
					if err := reportToRecorder(ctx, queueId, item.EventId, tenantId, code, result); err != nil {
						logger.Err(err).Msg("error reporting to recorder")
					}
					if err := makeDecisions(ctx, item.EventId, tenantId, code, result, nextSteps); err != nil {
						logger.Err(err).Msg("error enqueuing the result")
					}
				}
			}
		}
	}
}

func makeDecisions(ctx context.Context, eventId string, tenantId string, code uint64, result string, resultDef *pb.ResultDef) error {
	r := pb.JobResult{
		Code: code,
		//Result: result,
	}
	bytes1, err := proto.Marshal(&r)
	if err != nil {
		return err
	}
	var q *pb.QueueRequest
	if code == NO_ERROR {
		q = &pb.QueueRequest{
			TenantId: tenantId,
			QueueId:  resultDef.Ok.SupplierQueueId,
			Items: []*pb.QueueItem{{
				EventId: resultDef.Ok.ID,
				Data:    bytes1,
			},
			},
		}
	} else {
		q = &pb.QueueRequest{
			TenantId: tenantId,
			QueueId:  resultDef.Error.SupplierQueueId,
			Items: []*pb.QueueItem{{
				EventId: resultDef.Error.ID,
				Data:    bytes1,
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

func reportToRecorder(ctx context.Context, queueId string, eventId string, tenantId string, code uint64, result string) error {
	now := time.Now()
	host, err := os.Hostname()
	if err != nil {
		host = "<error>"
	}
	ex := &pb.JobExecution{
		EventId:  eventId,
		TenantId: tenantId,
		QueueId:  queueId,
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
