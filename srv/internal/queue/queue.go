package queue

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/utilico"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
)

type QueueStore[T any] struct {
	queues *utilico.SyncMap[string, Queue[T]]
}

type Queue[T any] interface {
	Add(data T) error
	Remove() (T, error)
}

func NewQueueStore[T any](ctx context.Context) (*QueueStore[T], error) {
	q := utilico.NewSyncMap[string, Queue[T]]()
	qs := &QueueStore[T]{
		queues: q,
	}
	if err := qs.load(ctx); err != nil {
		return nil, err
	}
	qs.startListeningUpdates(ctx)
	return qs, nil
}

func (q *QueueStore[T]) load(ctx context.Context) error {
	controlClient, err := remote.NewControlClient(ctx)
	if err != nil {
		return err
	}
	pkgs, err := controlClient.GetAllPackages(ctx)
	if err != nil {
		return err
	}
	q.addOrUpdateQueues(ctx, pkgs)
	if err != nil {
		return err
	}
	return nil
}

func (j *QueueStore[T]) startListeningUpdates(ctx context.Context) error {
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
func (j *QueueStore[T]) onUpdate(ctx context.Context, u *pb.UpdateToPackagesStrReply) {
	switch u.Type {
	case pb.UpdateType_Delete:
		j.deleteQueues(u.Object.TenantId, u.Object.Queues)
	case pb.UpdateType_New, pb.UpdateType_Update:
		j.addOrUpdateQueues(ctx, []*pb.JobPackage{u.Object})
	}
}
func (j *QueueStore[T]) deleteQueues(tenantId string, qs []*pb.QueueDef) {
	for _, q := range qs {
		j.queues.Delete(getFullQueueId(tenantId, q.ID))
	}
}

func (j *QueueStore[T]) addOrUpdateQueues(ctx context.Context, pkgs []*pb.JobPackage) error {
	for _, ps := range pkgs {
		tenantId := ps.TenantId
		for _, queue := range ps.Queues {
			err := j.addOrUpdateQueue(ctx, tenantId, queue)
			if err != nil {
				j.queues = utilico.NewSyncMap[string, Queue[T]]()
				return err
			}
		}
	}
	return nil
}

func (j *QueueStore[T]) addOrUpdateQueue(ctx context.Context, tenantId string, q *pb.QueueDef) error {
	queueId := getFullQueueId(tenantId, q.ID)
	queue, err := j.newQueue(queueId)
	if err != nil {
		return err
	}
	if j.existQueue(queueId) {
		j.queues.Swap(queueId, queue)
	} else {
		j.queues.Store(queueId, queue)
	}
	return nil
}

func (j *QueueStore[T]) newQueue(id string) (Queue[T], error) {
	return GetFileBasedQueue[T](id), nil
}

func getFullQueueId(tenantId string, queueId string) string {
	return tenantId + "/" + queueId
}

func (j *QueueStore[T]) existQueue(queueId string) bool {
	_, ok := j.queues.Load(queueId)
	return ok
}

var (
	ErrQueueUnknown = errors.New("queue unknown")
)

func (j *QueueStore[T]) GetQueue(tentantId string, queueId string) (Queue[T], error) {
	q, ok := j.queues.Load(getFullQueueId(tentantId, queueId))
	if !ok {
		return nil, ErrQueueUnknown
	}
	return q, nil
}
