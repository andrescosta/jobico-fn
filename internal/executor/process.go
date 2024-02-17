package executor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/runtimes/wasm"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

type process struct {
	packageID string
	tenant    string
	queue     string
	events    map[string]*event
	runtime   *wasm.Runtime
	cli       *cli
}

type event struct {
	id        string
	nextStep  *pb.ResultDef
	module    *module
	logSender *recorder
}

type module struct {
	id         uint32
	wasmModule *wasm.Module
}

func (p *process) processEvents(ctx context.Context, w *sync.WaitGroup) {
	defer w.Done()
	logger := zerolog.Ctx(ctx)
	items, err := p.cli.queue.Dequeue(ctx, p.tenant, p.queue)
	// TODO: do something with errors
	if err != nil || len(items) == 0 {
		return
	}
	for _, item := range items {
		event, ok := p.events[item.Event]
		if !ok {
			logger.Warn().Msgf("event %s not supported", event.id)
			continue
		}
		code, result, err := run(ctx, event.module.wasmModule, event.module.id, item.Data)
		if err != nil {
			logger.Err(err).Msg("error executing")
		}
		if err := event.logSender.sendResult(ctx, p.queue, code, result); err != nil {
			logger.Err(err).Msg("error reporting to recorder")
		}
		if err := p.makeDecisions(ctx, p.tenant, code, event.nextStep); err != nil {
			logger.Err(err).Msg("error enqueuing the result")
		}
	}
}

func (p *process) makeDecisions(ctx context.Context, tenant string, code uint64, resultDef *pb.ResultDef) error {
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
	if err := p.cli.queue.Queue(ctx, q); err != nil {
		return err
	}
	return nil
}

func run(ctx context.Context, module *wasm.Module, id uint32, data []byte) (uint64, string, error) {
	mod := "goenv"
	logger := zerolog.Ctx(ctx)
	ctx, cancel := context.WithTimeout(ctx, *env.Duration("wasm.exec.timeout", 2*time.Minute))
	defer cancel()
	code, result, err := module.Run(ctx, id, string(data))
	if err != nil {
		return 0, "", errors.Join(err, fmt.Errorf("error in module %s", mod))
	}
	logger.Debug().Msgf("%d | %s", code, result)
	return code, result, nil
}
