package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc/cache"
	evcache "github.com/andrescosta/goico/pkg/service/grpc/cache/event"
	"github.com/andrescosta/jobico/internal/api/remote"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/pkg/grpchelper"
)

type event struct {
	Data []interface{} `json:"data"`
}

type errSend struct {
	StatusCode int
}

type client struct {
	ctx          context.Context
	cacheEventCh chan *evcache.Event
	ctl          *remote.ControlClient
	repo         *remote.RepoClient
	queue        *remote.QueueClient
	recorder     *remote.RecorderClient
	cache        evcache.CacheServiceClient
	listener     *http.Client
}

var errTimeoutWaiting = errors.New("timeout while waiting channel")

func newClient(ctx context.Context, dialer service.GrpcDialer, transport service.HTTPTranporter) (*client, error) {
	ctl, err := remote.NewControlClient(ctx, dialer)
	if err != nil {
		return nil, err
	}
	queue, err := remote.NewQueueClient(ctx, dialer)
	if err != nil {
		return nil, err
	}
	repo, err := remote.NewRepoClient(ctx, dialer)
	if err != nil {
		return nil, err
	}
	recorder, err := remote.NewRecorderClient(ctx, dialer)
	if err != nil {
		return nil, err
	}
	cacheClient, err := cache.NewCacheServiceClient(ctx, env.String("cache_listener.addr"), dialer)
	if err != nil {
		return nil, err
	}
	listenerClient, err := newHttpClient(transport)
	if err != nil {
		return nil, err
	}
	return &client{
		ctx:      ctx,
		ctl:      ctl,
		queue:    queue,
		repo:     repo,
		recorder: recorder,
		cache:    cacheClient,
		listener: listenerClient,
	}, nil
}

func newHttpClient(t service.HTTPTranporter) (*http.Client, error) {
	addr := env.String("listener.host")
	transport, err := t.Tranport(addr)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Timeout:   1 * time.Second,
		Transport: transport,
	}, nil
}

type schemaRefIds struct {
	schameRef      string
	schemaRefOk    string
	schemaRefError string
}

func (s *client) close() error {
	errs := make([]error, 0)
	if s.cacheEventCh != nil {
		close(s.cacheEventCh)
	}
	if err := s.ctl.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.queue.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.recorder.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.repo.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (s *client) newTestPackage(schemaRefIds schemaRefIds, runtimeRef string) *pb.JobPackage {
	p := pb.JobPackage{
		ID:     "job_id_1",
		Name:   strptr("job_name_1"),
		Tenant: "tenant_1",
		Queues: []*pb.QueueDef{
			{
				ID:   "queue_id_1",
				Name: strptr("queue_name_1"),
			},
		},
		Jobs: []*pb.JobDef{
			{
				Event: &pb.EventDef{
					ID:            "event_id_1",
					Name:          strptr("event_name_1"),
					DataType:      pb.DataType_Json,
					SupplierQueue: "queue_id_1",
					Runtime:       "runtime_id_1",
					Schema:        &pb.SchemaDef{SchemaRef: schemaRefIds.schameRef},
				},
				Result: &pb.ResultDef{
					Ok: &pb.EventDef{
						ID: "event_id_1_ok",
					},
					Error: &pb.EventDef{
						ID: "event_id_1_error",
					},
				},
			},
			{
				Event: &pb.EventDef{
					ID:            "event_id_1_ok",
					Name:          strptr("event_name_1_ok"),
					DataType:      pb.DataType_Json,
					SupplierQueue: "queue_id_1",
					Runtime:       "runtime_id_1",
					Schema:        &pb.SchemaDef{SchemaRef: schemaRefIds.schemaRefOk},
				},
			},
			{
				Event: &pb.EventDef{
					ID:            "event_id_1_error",
					Name:          strptr("event_name_1_error"),
					DataType:      pb.DataType_Json,
					SupplierQueue: "queue_id_1",
					Runtime:       "runtime_id_1",
					Schema:        &pb.SchemaDef{SchemaRef: schemaRefIds.schemaRefError},
				},
			},
		},
		Runtimes: []*pb.RuntimeDef{
			{
				ID:           "runtime_id_1",
				Name:         strptr("runtime_name_1"),
				ModuleRef:    runtimeRef,
				MainFuncName: strptr("event"),
				Type:         pb.RuntimeType_Wasm10,
				Platform:     pb.Platform_TinyGO.Enum(),
			},
		},
	}
	return &p
}

func strptr(n string) *string {
	return &n
}

func (s *client) addPackage(p *pb.JobPackage) error {
	ps, err := s.ctl.GetPackage(s.ctx, p.Tenant, &p.ID)
	if err != nil {
		return err
	}
	if len(ps) >= 1 {
		return errors.New("too many packages with the same id")
	}

	if _, err := s.ctl.AddPackage(s.ctx, p); err != nil {
		return err
	}
	return nil
}

func (s *client) updatePackage(p *pb.JobPackage) error {
	ps, err := s.ctl.GetPackage(s.ctx, p.Tenant, &p.ID)
	if err != nil {
		return err
	}
	if len(ps) == 0 {
		return errors.New("package does not exist")
	}

	if err := s.ctl.UpdatePackage(s.ctx, p); err != nil {
		return err
	}
	return nil
}

func (s *client) deletePackage(p *pb.JobPackage) error {
	ps, err := s.ctl.GetPackage(s.ctx, p.Tenant, &p.ID)
	if err != nil {
		return err
	}
	if len(ps) == 0 {
		return errors.New("package does not exist")
	}

	if err := s.ctl.DeletePackage(s.ctx, p); err != nil {
		return err
	}
	return nil
}

func (s *client) uploadFile(tenant string, fileID string, fileType pb.File_FileType, f io.Reader) error {
	return s.repo.AddFile(s.ctx, tenant, fileID, fileType, f)
}

func (s *client) uploadSchemas(p *pb.JobPackage, files map[string][]byte) error {
	for _, e := range p.Jobs {
		schema := e.Event.Schema
		if err := s.uploadFile(p.Tenant,
			schema.SchemaRef,
			pb.File_JsonSchema, bytes.NewReader(files[schema.SchemaRef])); err != nil {
			return err
		}
	}
	return nil
}

func (s *client) sendEventV1(url *url.URL) error {
	d := struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Age       int    `json:"age"`
	}{"john", "connor", 50}
	ev := event{[]interface{}{d}}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return s.sendEvent(url, b)
}

func (s *client) sendEventV2(url *url.URL) error {
	d := struct {
		Name string `json:"name"`
		Dob  string `json:"dob"`
	}{"john connor", "01/01/1973"}
	ev := event{[]interface{}{d}}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return s.sendEvent(url, b)
}

func (s *client) sendEventMalFormed(url *url.URL) error {
	d := struct {
		Name string `json:"stname"`
		Age  string `json:"stage"`
	}{"john not connor", "01/01/2010"}
	ev := event{[]interface{}{d}}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return s.sendEvent(url, b)
}

func (e errSend) Error() string {
	return fmt.Sprintf("HTTP Status Code:%d", e.StatusCode)
}

func (s *client) sendEvent(url *url.URL, e []byte) error {
	r := &http.Request{
		Method: "POST",
		URL:    url,
		Body:   io.NopCloser(bytes.NewReader(e)),
	}
	re, err := s.listener.Do(r)
	if err != nil {
		return err
	}
	defer re.Body.Close()
	if re.StatusCode != http.StatusOK {
		return errSend{StatusCode: re.StatusCode}
	}
	return nil
}

func (s *client) addTenant(tenant string) error {
	t, err := s.ctl.GetTenant(s.ctx, &tenant)
	if err != nil {
		return err
	}
	if len(t) >= 1 {
		return errors.New("too many tenants with the same id")
	}
	_, err = s.ctl.AddTenant(s.ctx, &pb.Tenant{ID: tenant})
	if err != nil {
		return err
	}
	return nil
}

func (s *client) getAllPackages() ([]*pb.JobPackage, error) {
	t, err := s.ctl.GetAllPackages(s.ctx)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *client) dequeue(tenant string, queue string) ([]*pb.QueueItem, error) {
	res, err := s.queue.Dequeue(s.ctx, tenant, queue)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *client) startRecvCacheEvents() error {
	r, err := s.cache.Events(s.ctx, &evcache.Empty{})
	if err != nil {
		return err
	}
	ch := make(chan *evcache.Event)
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	go func() {
		defer func() {
			close(ch)
			cancel()
		}()
		_ = grpchelper.Recv[*evcache.Event](ctx, r, ch)
	}()
	s.cacheEventCh = ch
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (s *client) waitForCacheEvents() error {
	return waitFor(s.ctx, s.cacheEventCh)
}

func waitFor[T any](ctx context.Context, ch <-chan T) error {
	if ch == nil {
		return errors.New("channel is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		return errTimeoutWaiting
	case <-ch:
		return nil
	}
}
