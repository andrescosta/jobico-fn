package listener

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc/cache"
	"github.com/andrescosta/jobico/internal/api/client"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/rs/zerolog"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

var ErrEventUnknown = fmt.Errorf("event unknown")

type EventDefCache struct {
	eventCache    *cache.Cache[string, *EventEntry]
	serviceCache  *cache.Service
	repoClient    *client.Repo
	dialer        service.GrpcDialer
	controlClient *client.Ctl
}
type EventEntry struct {
	EventDef *pb.EventDef
	Schema   *jsonschema.Schema
}

func newCache(ctx context.Context, d service.GrpcDialer, l service.GrpcListener) (*EventDefCache, error) {
	repoClient, err := client.NewRepo(ctx, d)
	if err != nil {
		return nil, err
	}
	controlClient, err := client.NewCtl(ctx, d)
	if err != nil {
		return nil, err
	}
	b := cache.New[string, *EventEntry](ctx, "listener")
	svc, err := cache.NewService[string, *EventEntry](ctx, l, b)
	if err != nil {
		return nil, err
	}
	cache := EventDefCache{
		eventCache:    b,
		repoClient:    repoClient,
		controlClient: controlClient,
		dialer:        d,
		serviceCache:  svc,
	}
	if err := cache.populate(ctx); err != nil {
		return nil, err
	}
	if err := cache.startListeningUpdates(ctx, d); err != nil {
		return nil, err
	}
	go func() {
		logger := zerolog.Ctx(ctx)
		if err := cache.serviceCache.Serve(); err != nil {
			logger.Warn().Msgf("Error stopping cache: %v", err)
		}
	}()
	return &cache, nil
}

func (j *EventDefCache) Close() error {
	var err error
	err = errors.Join(err, j.controlClient.Close())
	err = errors.Join(err, j.repoClient.Close())
	err = errors.Join(err, j.eventCache.Close())
	return err
}

func (j *EventDefCache) Get(tenant string, eventID string) (*EventEntry, error) {
	ev, ok := j.eventCache.Get(getEventName(tenant, eventID))
	if !ok {
		return nil, ErrEventUnknown
	}
	return ev, nil
}

func (j *EventDefCache) startListeningUpdates(ctx context.Context, d service.GrpcDialer) error {
	l, err := j.controlClient.ListenerForPackageUpdates(ctx)
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

func (j *EventDefCache) onUpdate(ctx context.Context, u *pb.UpdateToPackagesStrReply) {
	switch u.Type {
	case pb.UpdateType_Delete:
		j.deleteEventsOfPackage(ctx, u.Object)
	case pb.UpdateType_New, pb.UpdateType_Update:
		j.addOrUpdateEventsForPackages(ctx, []*pb.JobPackage{u.Object})
	}
}

func (j *EventDefCache) deleteEventsOfPackage(ctx context.Context, p *pb.JobPackage) {
	logger := zerolog.Ctx(ctx)
	for _, d := range p.Jobs {
		if err := j.eventCache.Delete(getEventName(p.Tenant, d.Event.ID)); err != nil {
			logger.Warn().Msgf("cache broken, error: %v", err)
		}
	}
}

func (j *EventDefCache) populate(ctx context.Context) error {
	pkgs, err := j.controlClient.AllPackages(ctx)
	if err != nil {
		return err
	}
	return j.addOrUpdateEventsForPackages(ctx, pkgs)
}

func (j *EventDefCache) addOrUpdateEventsForPackages(ctx context.Context, pkgs []*pb.JobPackage) error {
	for _, ps := range pkgs {
		tenant := ps.Tenant
		for _, job := range ps.Jobs {
			if err := j.addOrUpdateEvent(ctx, tenant, job); err != nil {
				return err
			}
		}
	}
	return nil
}

func (j *EventDefCache) addOrUpdateEvent(ctx context.Context, tenant string, job *pb.JobDef) error {
	logger := zerolog.Ctx(ctx)
	event := job.Event
	f, err := j.repoClient.File(ctx, tenant, event.Schema.SchemaRef)
	if err != nil {
		return err
	}
	comp := jsonschema.NewCompiler()
	if err := comp.AddResource(getEventName(tenant, event.ID), bytes.NewReader(f)); err != nil {
		return err
	}
	compiledSchema, err := comp.Compile(getEventName(tenant, event.ID))
	if err != nil {
		return err
	}
	eventName := getEventName(tenant, event.ID)
	ev := &EventEntry{
		EventDef: event,
		Schema:   compiledSchema,
	}

	if err := j.eventCache.AddOrUpdate(eventName, ev); err != nil {
		logger.Warn().Msgf("cache broken, error: %v", err)
	}
	return nil
}

func getEventName(tenant string, eventID string) string {
	return tenant + "/" + eventID
}
