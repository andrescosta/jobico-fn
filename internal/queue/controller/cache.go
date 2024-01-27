package controller

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/collection"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/jobico/internal/api/remote"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/queue/provider"
	"github.com/rs/zerolog"
)

var ErrQueueUnknown = errors.New("queue unknown")

type Option struct {
	InMemory bool
}

type QueueBuilder[T any] func(string) provider.Queue[T]

// TODO: fix locks!! cal LoadOrSTore or what ?

type Cache[T any] struct {
	queues        *collection.SyncMap[string, provider.Queue[T]]
	dialer        service.GrpcDialer
	newQueue      QueueBuilder[T]
	controlClient *remote.CtlClient
}

func NewQueueCache[T any](ctx context.Context, d service.GrpcDialer, o Option) (*Cache[T], error) {
	q := collection.NewSyncMap[string, provider.Queue[T]]()
	b := func(id string) provider.Queue[T] { return provider.NewFileQueue[T](id) }
	if o.InMemory {
		b = func(_ string) provider.Queue[T] { return provider.NewMemBasedQueue[T]() }
	}
	c, err := remote.NewCtlClient(ctx, d)
	if err != nil {
		return nil, err
	}
	qs := &Cache[T]{
		queues:        q,
		dialer:        d,
		newQueue:      b,
		controlClient: c,
	}
	if err := qs.populate(ctx); err != nil {
		return nil, err
	}
	if err := qs.startListeningUpdates(ctx); err != nil {
		return nil, err
	}
	return qs, nil
}

func (q *Cache[T]) Close() error {
	return q.controlClient.Close()
}

func (q *Cache[T]) GetQueue(tentant string, queueID string) (provider.Queue[T], error) {
	queue, ok := q.queues.Load(getQueueName(tentant, queueID))
	if !ok {
		return nil, ErrQueueUnknown
	}
	return queue, nil
}

func (q *Cache[T]) populate(ctx context.Context) error {
	pkgs, err := q.controlClient.GetAllPackages(ctx)
	if err != nil {
		return err
	}
	if err := q.addQueues(ctx, pkgs); err != nil {
		return err
	}
	return nil
}

func (q *Cache[T]) startListeningUpdates(ctx context.Context) error {
	l, err := q.controlClient.ListenerForPackageUpdates(ctx)
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

func (q *Cache[T]) onUpdate(ctx context.Context, u *pb.UpdateToPackagesStrReply) {
	logger := zerolog.Ctx(ctx)
	switch u.Type {
	case pb.UpdateType_Delete:
		q.deleteQueues(u.Object.Tenant, u.Object.Queues)
	case pb.UpdateType_New, pb.UpdateType_Update:
		if err := q.addQueues(ctx, []*pb.JobPackage{u.Object}); err != nil {
			logger.Warn().AnErr("error", err).Msg("onUpdate: error updating queue")
		}
	}
}

func (q *Cache[T]) deleteQueues(tenant string, qs []*pb.QueueDef) {
	for _, queue := range qs {
		q.queues.Delete(getQueueName(tenant, queue.ID))
	}
}

func (q *Cache[T]) addQueues(ctx context.Context, pkgs []*pb.JobPackage) error {
	for _, ps := range pkgs {
		tenant := ps.Tenant
		for _, queue := range ps.Queues {
			q.addQueue(ctx, tenant, queue)
		}
	}
	return nil
}

func (q *Cache[T]) addQueue(_ context.Context, tenant string, def *pb.QueueDef) {
	name := getQueueName(tenant, def.ID)
	queue := q.newQueue(name)
	_ = q.queues.LoadOrStore(name, queue)
}

func getQueueName(tenant string, queueID string) string {
	return tenant + "/" + queueID
}
