package listener

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/rs/zerolog"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type EventDefCache struct {
	defs       *sync.Map
	repoClient *remote.RepoClient
}
type EventEntry struct {
	EventDef *pb.EventDef
	Schema   *jsonschema.Schema
}

func NewEventDefCache(ctx context.Context) (*EventDefCache, error) {
	repoClient, err := remote.NewRepoClient(ctx)
	if err != nil {
		return nil, err
	}
	store := EventDefCache{
		defs:       &sync.Map{},
		repoClient: repoClient,
	}
	if err := store.fill(ctx); err != nil {
		return nil, err
	}
	if err := store.startListeningUpdates(ctx); err != nil {
		return nil, err
	}
	return &store, nil
}

func (j *EventDefCache) Get(tenant string, eventID string) (*EventEntry, error) {
	ev, ok := j.defs.Load(getEventName(tenant, eventID))
	if !ok {
		return nil, ErrEventUnknown
	}
	res, ok := ev.(*EventEntry)
	if !ok {
		return nil, ErrEventUnknown
	}
	return res, nil
}
func (j *EventDefCache) startListeningUpdates(ctx context.Context) error {
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

func (j *EventDefCache) onUpdate(ctx context.Context, u *pb.UpdateToPackagesStrReply) {
	switch u.Type {
	case pb.UpdateType_Delete:
		j.deleteEventsOfPackage(u.Object)
	case pb.UpdateType_New, pb.UpdateType_Update:
		j.addOrUpdateEventsForPackages(ctx, []*pb.JobPackage{u.Object})
	}
}

func (j *EventDefCache) deleteEventsOfPackage(p *pb.JobPackage) {
	for _, d := range p.Jobs {
		j.defs.Delete(getEventName(p.Tenant, d.Event.ID))
	}
}

func (j *EventDefCache) fill(ctx context.Context) error {
	controlClient, err := remote.NewControlClient(ctx)
	if err != nil {
		return err
	}
	pkgs, err := controlClient.GetAllPackages(ctx)
	if err != nil {
		return err
	}
	j.addOrUpdateEventsForPackages(ctx, pkgs)
	return nil
}

func (j *EventDefCache) addOrUpdateEventsForPackages(ctx context.Context, pkgs []*pb.JobPackage) {
	logger := zerolog.Ctx(ctx)
	for _, ps := range pkgs {
		tenant := ps.Tenant
		for _, job := range ps.Jobs {
			if err := j.addOrUpdateEvent(ctx, tenant, job); err != nil {
				logger.Warn().AnErr("error", err).Msg("Error updating package")
			}
		}
	}
}

func (j *EventDefCache) addOrUpdateEvent(ctx context.Context, tenant string, job *pb.JobDef) error {
	event := job.Event
	f, err := j.repoClient.GetFile(ctx, tenant, event.Schema.SchemaRef)
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
	if j.existEvent(eventName) {
		j.defs.Swap(eventName, ev)
	} else {
		j.defs.Store(
			eventName,
			ev)
	}
	return nil
}

var (
	ErrEventUnknown = fmt.Errorf("event unknown")
)

func (j *EventDefCache) existEvent(eventName string) bool {
	_, ok := j.defs.Load(eventName)
	return ok
}
func getEventName(tenant string, eventID string) string {
	return tenant + "/" + eventID
}
