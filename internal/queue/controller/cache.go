package controller

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/collection"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/queue/provider"
	"github.com/rs/zerolog"
)

var (
	ErrQueueUnknown = errors.New("queue unknown")
)

type QueueCache[T any] struct {
	queues *collection.SyncMap[string, Queue[T]]
}
type Queue[T any] interface {
	Add(data T) error
	Remove() (T, error)
}

func NewQueueCache[T any](ctx context.Context) (*QueueCache[T], error) {
	q := collection.NewSyncMap[string, Queue[T]]()
	qs := &QueueCache[T]{
		queues: q,
	}
	if err := qs.load(ctx); err != nil {
		return nil, err
	}
	if err := qs.startListeningUpdates(ctx); err != nil {
		return nil, err
	}
	return qs, nil
}
func (q *QueueCache[T]) GetQueue(tentant string, queueID string) (Queue[T], error) {
	queue, ok := q.queues.Load(getQueueName(tentant, queueID))
	if !ok {
		return nil, ErrQueueUnknown
	}
	return queue, nil
}
func (q *QueueCache[T]) load(ctx context.Context) error {
	controlClient, err := remote.NewControlClient(ctx)
	if err != nil {
		return err
	}
	pkgs, err := controlClient.GetAllPackages(ctx)
	if err != nil {
		return err
	}
	if err := q.addOrUpdateQueuesForJobPackge(ctx, pkgs); err != nil {
		return err
	}
	return nil
}
func (q *QueueCache[T]) startListeningUpdates(ctx context.Context) error {
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
				q.onUpdate(ctx, u)
			}
		}
	}()
	return nil
}
func (q *QueueCache[T]) onUpdate(ctx context.Context, u *pb.UpdateToPackagesStrReply) {
	logger := zerolog.Ctx(ctx)
	switch u.Type {
	case pb.UpdateType_Delete:
		q.deleteQueues(u.Object.Tenant, u.Object.Queues)
	case pb.UpdateType_New, pb.UpdateType_Update:
		if err := q.addOrUpdateQueuesForJobPackge(ctx, []*pb.JobPackage{u.Object}); err != nil {
			logger.Warn().AnErr("error", err).Msg("onUpdate: error updating queue")
		}
	}
}
func (q *QueueCache[T]) deleteQueues(tenant string, qs []*pb.QueueDef) {
	for _, queue := range qs {
		q.queues.Delete(getQueueName(tenant, queue.ID))
	}
}
func (q *QueueCache[T]) addOrUpdateQueuesForJobPackge(ctx context.Context, pkgs []*pb.JobPackage) error {
	for _, ps := range pkgs {
		tenant := ps.Tenant
		for _, queue := range ps.Queues {
			err := q.addOrUpdateQueue(ctx, tenant, queue)
			if err != nil {
				q.queues = collection.NewSyncMap[string, Queue[T]]()
				return err
			}
		}
	}
	return nil
}
func (q *QueueCache[T]) addOrUpdateQueue(_ context.Context, tenant string, def *pb.QueueDef) error {
	name := getQueueName(tenant, def.ID)
	queue, err := q.newQueue(name)
	if err != nil {
		return err
	}
	if q.existQueue(name) {
		q.queues.Swap(name, queue)
	} else {
		q.queues.Store(name, queue)
	}
	return nil
}
func (q *QueueCache[T]) newQueue(id string) (Queue[T], error) {
	return provider.GetFileBasedQueue[T](id), nil
}
func (q *QueueCache[T]) existQueue(queueID string) bool {
	_, ok := q.queues.Load(queueID)
	return ok
}

func getQueueName(tenant string, queueID string) string {
	return tenant + "/" + queueID
}
