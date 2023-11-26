package executor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/workflew/api/pkg/remote"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/srv/internal/wasi"
	"github.com/rs/zerolog"
)

const (
	cacheDir = ".\\cache"
)

type jobPackage struct {
	PackageId string
	TenantId  string
	Queues    []string
	Runtime   *wasi.WasmRuntime
	Modules   map[string]*wasi.WasmModuleString
}

func getPackages(ctx context.Context) ([]*jobPackage, error) {
	ps, err := remote.NewControlClient().GetAllPackages(ctx)
	if err != nil {
		return nil, err
	}
	jobPackages := make([]*jobPackage, 0)
	for _, pkg := range ps {
		jobPackage := &jobPackage{}
		jobPackage.TenantId = pkg.TenantId
		jobPackage.PackageId = pkg.JobPackageId
		jobPackages = append(jobPackages, jobPackage)

		queues := make([]string, 0)
		modulesForEvents := make(map[string]*wasi.WasmModuleString)

		jobPackage.Runtime, err = wasi.NewWasmRuntime(ctx, cacheDir)
		if err != nil {
			return nil, err
		}
		for _, q := range pkg.Queues {
			queues = append(queues, q.QueueId)
		}
		jobPackage.Queues = queues

		repoClient := remote.NewRepoClient()
		files := make(map[string][]byte)
		for _, job := range pkg.Jobs {
			event := job.Event
			for _, runtime := range pkg.Runtimes {
				if runtime.RuntimeId == event.RuntimeId {
					wasmfile, ok := files[runtime.ModuleRef]
					if !ok {
						wasmfile, err = repoClient.GetFile(ctx, pkg.TenantId, runtime.ModuleRef)
						if err != nil {
							return nil, err
						}
						files[runtime.ModuleRef] = wasmfile
					}
					// We create a module per queue because wazero module call is not goroutine compatible
					wasmModule, err := wasi.NewWasmModuleString(ctx, jobPackage.Runtime, wasmfile, runtime.MainFuncName)
					if err != nil {
						return nil, err
					}
					modulesForEvents[getModuleName(event.SupplierQueueId, event.EventId)] = wasmModule
					break
				}
			}
		}
		jobPackage.Modules = modulesForEvents
	}
	return jobPackages, nil
}

func getModuleName(supplierQueueId string, eventId string) string {
	return supplierQueueId + "/" + eventId

}
func StartExecutors(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	pkgs, err := getPackages(ctx)
	if err != nil {
		return err
	}
	defer func() {
		for _, pkg := range pkgs {
			for _, m := range pkg.Modules {
				m.Close(ctx)
			}
			pkg.Runtime.Close(ctx)
		}
	}()

	var w sync.WaitGroup
	for _, pkg := range pkgs {
		for _, queue := range pkg.Queues {
			w.Add(1)
			go executor(ctx, pkg.TenantId, queue, pkg.Modules, &w)
		}
	}
	logger.Info().Msg("Workers started")
	w.Wait()
	logger.Info().Msg("Workers stoped")
	return nil
}

func executor(ctx context.Context, tenantId string, queueId string, modules map[string]*wasi.WasmModuleString, w *sync.WaitGroup) {
	defer w.Done()
	logger := zerolog.Ctx(ctx)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	logger.Debug().Msgf("Worker for tenant: %s and queue: %s started", tenantId, queueId)
	queueErrors := 0
	maxQueueErrors := env.GetAsInt("max.queue.errors", 10)
	for {
		select {
		case <-ctx.Done():
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
				if execute(ctx, module, item.Data); err != nil {
					logger.Err(err).Msg("error executing")
				}
			}
		}
	}
}

func dequeue(ctx context.Context, tenant string, queue string) ([]*pb.QueueItem, error) {
	return remote.NewQueueClient().Dequeue(ctx, tenant, queue)
}

func execute(ctx context.Context, module *wasi.WasmModuleString, data []byte) error {
	mod := "goenv"
	logger := zerolog.Ctx(ctx)
	out, err := module.ExecuteMainFunc(ctx, string(data))
	if err != nil {
		return errors.Join(err, fmt.Errorf("error in module %s", mod))
	}
	logger.Debug().Msg(out)
	return nil
}
