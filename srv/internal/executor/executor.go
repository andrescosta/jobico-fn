package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/andrescosta/goico/pkg/io"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/srv/internal/wasi"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Queue struct {
	MerchantId string
	QueueId    string
	Func       string
}

func getQueues(path string) ([]*Queue, error) {
	mechants, err := io.GetSubDirectories(path)
	queues := make([]*Queue, 0)
	if err != nil {
		return nil, err
	}
	for _, merchant := range mechants {
		qs, _ := io.GetSubDirectories(io.BuildFullPath([]string{path, merchant}))
		for _, q := range qs {
			queues = append(queues, &Queue{
				MerchantId: merchant,
				QueueId:    q,
				Func:       "greet",
			})
		}
	}
	return queues, nil
}

func StartExecutors(ctx context.Context, path string) error {
	logger := zerolog.Ctx(ctx)
	qs, err := getQueues(path)
	if err != nil {
		return err
	}
	var w sync.WaitGroup
	runtime, err := wasi.NewWasmRuntime(ctx, ".\\", []string{"greet"})
	defer runtime.Close(ctx)
	if err != nil {
		return err
	}
	for _, q := range qs {
		w.Add(1)
		go executor(ctx, q.MerchantId, q.QueueId, q.Func, runtime, &w)
	}
	logger.Info().Msg("Workers started")
	w.Wait()
	logger.Info().Msg("Workers stoped")
	return nil
}

func executor(ctx context.Context, merchant string, queue string, wasmFunc string, runtime *wasi.WasmRuntime, w *sync.WaitGroup) {
	defer w.Done()
	logger := zerolog.Ctx(ctx)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	logger.Debug().Msgf("Worker for merchant: %s and queue: %s started", merchant, queue)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d, err := query(ctx, merchant, queue)
			if err != nil {
				logger.Err(err)
			} else {
				for _, ds := range d {
					if err = execute(ctx, ds, wasmFunc, runtime); err != nil {
						logger.Debug().Msg(err.Error())
						logger.Err(err)
					}
				}
			}
		}
	}

}

func query(ctx context.Context, merchant string, queue string) ([]string, error) {
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
		QueueId: &pb.QueueId{
			Name: queue,
		},
		MerchantId: &pb.MerchantId{
			Id: merchant,
		},
	}

	r, err := client.Dequeue(ctx, &request)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, it := range r.Items {
		ret = append(ret, it.Data)
	}
	return ret, nil
}

func execute(ctx context.Context, data string, wasmFunc string, runtime *wasi.WasmRuntime) error {
	mod := "goenv"
	logger := zerolog.Ctx(ctx)
	out, err := runtime.Event(ctx, wasmFunc, data) //wasi.InvokeModule(ctx, mod, data)
	if err != nil {
		return errors.Join(err, fmt.Errorf("error in module %s", mod))
	}
	logger.Debug().Msg(out)
	return nil
}