package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/andrescosta/workflew/api/pkg/remote"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/srv/internal/wasi"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type executionQueue struct {
	TenantId      string
	QueueId       string
	FuncPerEvents map[string]*exFunc
}
type exFunc struct {
	Name     string
	ModuleId string
	Runtime  *wasi.WasmRuntime
}

func getRuntime(runtimeId string, rds []*pb.RuntimeDef) *pb.RuntimeDef {
	for _, r := range rds {
		if runtimeId == r.RuntimeId {
			return r
		}
	}
	return nil
}

func getPackages(ctx context.Context) ([]*executionQueue, error) {
	ps, err := remote.NewControlClient().GetAllPackages(ctx)
	if err != nil {
		return nil, err
	}
	pkgs := make([]*executionQueue, 0)
	repoClient := remote.NewRepoClient()
	for _, p := range ps {
		exs := make(map[string]*exFunc)
		for _, ex := range p.Executors {
			runtime := getRuntime(ex.RuntimeId, p.Runtimes)

			wasmfile, err := repoClient.GetFile(ctx, p.TenantId, runtime.ModuleRef)
			if err != nil {
				return nil, err
			}

			// TODO: we hsould have only one runtime type
			wa := []*wasi.Func{
				{
					ModuleId:   runtime.RuntimeId,
					WasmModule: wasmfile,
					FuncName:   ex.FuncName,
				},
			}

			wruntime, err := wasi.NewWasmRuntime(ctx, ".\\", wa)
			if err != nil {
				return nil, err
			}
			exf := &exFunc{
				Name:     ex.FuncName,
				Runtime:  wruntime,
				ModuleId: runtime.RuntimeId,
			}
			for _, e := range ex.SupportedEvents {
				exs[e] = exf
			}
		}
		for _, qq := range p.Queues {
			pkgs = append(pkgs,
				&executionQueue{
					TenantId:      p.TenantId,
					QueueId:       qq.QueueId,
					FuncPerEvents: exs,
				})
		}
	}
	return pkgs, nil
}

func StartExecutors(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	exqs, err := getPackages(ctx)
	if err != nil {
		return err
	}
	var w sync.WaitGroup

	//defer runtime.Close(ctx)
	if err != nil {
		return err
	}
	for _, q := range exqs {
		w.Add(1)
		go executor(ctx, q.TenantId, q.QueueId, q.FuncPerEvents, &w)
	}
	logger.Info().Msg("Workers started")
	w.Wait()
	logger.Info().Msg("Workers stoped")
	return nil
}

func executor(ctx context.Context, tenantId string, queueId string, exs map[string]*exFunc, w *sync.WaitGroup) {
	defer w.Done()
	logger := zerolog.Ctx(ctx)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	logger.Debug().Msgf("Worker for tenant: %s and queue: %s started", tenantId, queueId)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d, err := query(ctx, tenantId, queueId)
			if err != nil {
				logger.Err(err).Msg("Error quering")
			} else {
				for _, ds := range d {
					runtime := exs[ds.EventId]
					if err = execute(ctx, runtime.ModuleId, ds.Data, runtime.Runtime); err != nil {
						logger.Debug().Msg(err.Error())
						logger.Err(err).Msg("error executing")
					}
				}
			}
		}
	}

}

func query(ctx context.Context, tenant string, queue string) ([]*pb.QueueItem, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	Host := os.Getenv("queue.host")
	conn, err := grpc.Dial(Host, opts...)

	if err != nil {
		return nil, err
	}
	client := pb.NewQueueClient(conn)

	defer conn.Close()

	request := pb.DequeueRequest{
		QueueId:  queue,
		TenantId: tenant,
	}

	r, err := client.Dequeue(ctx, &request)
	if err != nil {
		return nil, err
	}
	return r.Items, nil
}

func execute(ctx context.Context, wasmRuntime string, data []byte, runtime *wasi.WasmRuntime) error {
	mod := "goenv"
	logger := zerolog.Ctx(ctx)
	out, err := runtime.Execute(ctx, wasmRuntime, data) //wasi.InvokeModule(ctx, mod, data)
	if err != nil {
		return errors.Join(err, fmt.Errorf("error in module %s", mod))
	}
	logger.Debug().Msg(out)
	return nil
}
