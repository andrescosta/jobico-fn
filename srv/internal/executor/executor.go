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

		for _, q := range pkg.Queues {
			queues = append(queues, q.QueueId)
		}
		jobPackage.Queues = queues

		jobPackage.Runtime, err = wasi.NewWasmRuntime(ctx, cacheDir)
		if err != nil {
			return nil, err
		}

		repoClient := remote.NewRepoClient()
		files := make(map[string][]byte)
		for _, executor := range pkg.Executors {
			var wasmModule *wasi.WasmModuleString
			for _, runtime := range pkg.Runtimes {
				if runtime.RuntimeId == executor.RuntimeId {
					wasmfile, ok := files[runtime.ModuleRef]
					if !ok {
						wasmfile, err = repoClient.GetFile(ctx, pkg.TenantId, runtime.ModuleRef)
						if err != nil {
							return nil, err
						}
						files[runtime.ModuleRef] = wasmfile
					}
					wasmModule, err = wasi.NewWasmModuleString(ctx, jobPackage.Runtime, wasmfile, runtime.MainFuncName)
					if err != nil {
						return nil, err
					}
					break
				}
			}
			for _, sevent := range executor.SupportedEvents {
				modulesForEvents[sevent] = wasmModule
			}
		}

		jobPackage.Modules = modulesForEvents
	}
	return jobPackages, nil
}

func StartExecutors(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	pkgs, err := getPackages(ctx)
	if err != nil {
		return err
	}
	defer func() {
		for _, pkg := range pkgs {
			pkg.Runtime.Close(ctx)
		}
	}()

	var w sync.WaitGroup
	for _, pkg := range pkgs {
		w.Add(1)
		for _, queue := range pkg.Queues {
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
				module, ok1 := modules[item.EventId]
				if !ok1 {
					logger.Warn().Msgf("event %s not supported", item.EventId)
					continue
				}
				if execute(ctx, module, item.Data); err != nil {
					logger.Debug().Msg(err.Error())
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
