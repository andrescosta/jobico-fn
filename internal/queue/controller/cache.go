package controller

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/collection"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/syncutil"
	"github.com/andrescosta/jobico/internal/api/client"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/queue/provider"
	"github.com/rs/zerolog"
)

var ErrQueueUnknown = errors.New("queue unknown")

type Option struct {
	InMemory bool
	Dir      string
}

type QueueBuilder[T any] func(string) (provider.Queue[T], error)

type Cache[T any] struct {
	init         *syncutil.OnceDisposable
	queues       *collection.SyncMap[string, provider.Queue[T]]
	queueBuilder QueueBuilder[T]
	ctl          *client.Ctl
}

func NewCache[T any](ctx context.Context, dialer service.GrpcDialer, o Option) (*Cache[T], error) {
	syncmap := collection.NewSyncMap[string, provider.Queue[T]]()
	queueBuilder := func(id string) (provider.Queue[T], error) { return provider.NewFileQueue[T](o.Dir, id) }
	if o.InMemory {
		queueBuilder = func(_ string) (provider.Queue[T], error) { return provider.NewMemBasedQueue[T]() }
	}
	ctl, err := client.NewCtl(ctx, dialer)
	if err != nil {
		return nil, err
	}
	cache := &Cache[T]{
		init:         syncutil.NewOnceDisposable(),
		queues:       syncmap,
		ctl:          ctl,
		queueBuilder: queueBuilder,
	}
	return cache, nil
}

func (q *Cache[T]) populate(ctx context.Context) error {
	if err := q.addPackages(ctx); err != nil {
		return err
	}
	if err := q.startListeningUpdates(ctx); err != nil {
		return err
	}
	return nil
}

func (q *Cache[T]) Close() error {
	err := q.init.Dispose(context.Background(), func(_ context.Context) error {
		if q.ctl != nil {
			return q.ctl.Close()
		}
		return nil
	})
	if errors.Is(err, syncutil.ErrTaskNotDone) {
		return nil
	}
	return err
}

func (q *Cache[T]) GetQueue(ctx context.Context, tentant string, queueID string) (provider.Queue[T], error) {
	err := q.init.Do(ctx, q.populate)
	if err != nil {
		return nil, err
	}
	queue, ok := q.queues.Load(getQueueName(tentant, queueID))
	if !ok {
		return nil, ErrQueueUnknown
	}
	return queue, nil
}

func (q *Cache[T]) addPackages(ctx context.Context) error {
	pkgs, err := q.ctl.AllPackages(ctx)
	if err != nil {
		return err
	}
	if err := q.addOrUpdate(ctx, pkgs); err != nil {
		return err
	}
	return nil
}

func (q *Cache[T]) startListeningUpdates(ctx context.Context) error {
	l, err := q.ctl.ListenerForPackageUpdates(ctx)
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
		q.delete(u.Object.Tenant, u.Object.Queues)
	case pb.UpdateType_New, pb.UpdateType_Update:
		if err := q.addOrUpdate(ctx, []*pb.JobPackage{u.Object}); err != nil {
			logger.Warn().AnErr("error", err).Msg("onUpdate: error updating queue")
		}
	}
}

func (q *Cache[T]) delete(tenant string, qs []*pb.QueueDef) {
	for _, queue := range qs {
		q.queues.Delete(getQueueName(tenant, queue.ID))
	}
}

func (q *Cache[T]) addOrUpdate(ctx context.Context, pkgs []*pb.JobPackage) error {
	for _, ps := range pkgs {
		tenant := ps.Tenant
		for _, queue := range ps.Queues {
			err := q.addOrUpdateQueue(ctx, tenant, queue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (q *Cache[T]) addOrUpdateQueue(ctx context.Context, tenant string, def *pb.QueueDef) error {
	name := getQueueName(tenant, def.ID)
	queue, err := q.queueBuilder(name)
	if err != nil {
		return err
	}
	zerolog.Ctx(ctx).Debug().Msgf("New queue:%s", def.ID)
	_ = q.queues.LoadOrStore(name, queue)
	return nil
}

func getQueueName(tenant string, queueID string) string {
	return tenant + "/" + queueID
}
