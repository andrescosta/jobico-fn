package listener

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type EventsStore struct {
	events     *sync.Map
	repoClient *remote.RepoClient
}

type Event struct {
	event  *pb.EventDef
	schema *jsonschema.Schema
}

func NewEventsStore(ctx context.Context) (*EventsStore, error) {
	repoClient, err := remote.NewRepoClient(ctx)
	if err != nil {
		return nil, err
	}
	store := EventsStore{
		events:     &sync.Map{},
		repoClient: repoClient,
	}

	if err := store.load(ctx); err != nil {
		return nil, err
	}

	store.startListeningUpdates(ctx)
	return &store, nil
}

func (j *EventsStore) startListeningUpdates(ctx context.Context) error {
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

func (j *EventsStore) onUpdate(ctx context.Context, u *pb.UpdateToPackagesStrReply) {
	switch u.Type {
	case pb.UpdateType_Delete:
		j.deletePackage(u.Object)
	case pb.UpdateType_New, pb.UpdateType_Update:
		j.addOrUpdatePackages(ctx, []*pb.JobPackage{u.Object})
	}
}
func (j *EventsStore) deletePackage(p *pb.JobPackage) {
	for _, d := range p.Jobs {
		j.events.Delete(getFullEventId(p.TenantId, d.Event.ID))
	}
}
func (j *EventsStore) load(ctx context.Context) error {
	controlClient, err := remote.NewControlClient(ctx)
	if err != nil {
		return err
	}
	pkgs, err := controlClient.GetAllPackages(ctx)
	if err != nil {
		return err
	}
	j.addOrUpdatePackages(ctx, pkgs)
	return nil
}
func (j *EventsStore) addOrUpdatePackages(ctx context.Context, pkgs []*pb.JobPackage) {
	for _, ps := range pkgs {
		tenantId := ps.TenantId
		for _, job := range ps.Jobs {
			j.addOrUpdateEvent(ctx, tenantId, job)
		}
	}
}

func (j *EventsStore) addOrUpdateEvent(ctx context.Context, tenantId string, job *pb.JobDef) error {
	event := job.Event
	f, err := j.repoClient.GetFile(ctx, tenantId, event.Schema.SchemaRef)
	if err != nil {
		return err
	}
	comp := jsonschema.NewCompiler()
	if err := comp.AddResource(getFullEventId(tenantId, event.ID), bytes.NewReader(f)); err != nil {
		return err
	}
	compiledSchema, err := comp.Compile(getFullEventId(tenantId, event.ID))
	if err != nil {
		return err
	}

	fullEventId := getFullEventId(tenantId, event.ID)
	ev := &Event{
		event:  event,
		schema: compiledSchema,
	}
	if j.existEvent(fullEventId) {
		j.events.Swap(fullEventId, ev)
	} else {
		j.events.Store(
			fullEventId,
			ev)
	}
	return nil
}

var (
	ErrEventUnknown = fmt.Errorf("event unknown")
)

func (j *EventsStore) GetEvent(tenantId string, eventId string) (*Event, error) {
	ev, ok := j.events.Load(getFullEventId(tenantId, eventId))
	if !ok {
		return nil, ErrEventUnknown
	}
	res, ok := ev.(*Event)
	if !ok {
		return nil, ErrEventUnknown
	}
	return res, nil
}

func (j *EventsStore) existEvent(fullEventId string) bool {
	_, ok := j.events.Load(fullEventId)
	return ok
}

func getFullEventId(tenantId string, eventId string) string {
	return tenantId + "/" + eventId
}
