package executor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/internal/utils"
	"github.com/andrescosta/workflew/internal/wazero"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Queue struct {
	MerchantId string
	QueueId    string
}

func getQueues(path string) ([]*Queue, error) {
	mechants, err := utils.GetSubDirectories(path)
	queues := make([]*Queue, 0)
	if err != nil {
		return nil, err
	}
	for _, merchant := range mechants {
		qs, _ := utils.GetSubDirectories(utils.BuildFullPath([]string{path, merchant}))
		for _, q := range qs {
			queues = append(queues, &Queue{
				MerchantId: merchant,
				QueueId:    q,
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
	for _, q := range qs {
		w.Add(1)
		go executor(ctx, q.MerchantId, q.QueueId, &w)
	}
	logger.Info().Msg("Workers started")
	w.Wait()
	logger.Info().Msg("Workers stoped")
	return nil
}

func executor(ctx context.Context, merchant string, queue string, w *sync.WaitGroup) {
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
					if err = execute(ctx, ds); err != nil {
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

func execute(ctx context.Context, data string) error {
	mod := "goenv"
	modpath := fmt.Sprintf("target/%v.wasm", mod)
	log.Printf("loading module %v", modpath)
	var env map[string]string
	env = make(map[string]string)
	env["data"] = data
	out, err := wazero.InvokeWasmModule(ctx, mod, modpath, env)
	if err != nil {

		return errors.Join(err, fmt.Errorf("error loading module %s", modpath))
	}
	fmt.Println(out)
	return nil
}
