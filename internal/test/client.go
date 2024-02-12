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
	"strings"
	"time"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc/cache"
	evcache "github.com/andrescosta/goico/pkg/service/grpc/cache/event"
	"github.com/andrescosta/goico/pkg/test"
	"github.com/andrescosta/jobico/internal/api/client"
	pb "github.com/andrescosta/jobico/internal/api/types"
)

type event struct {
	Data []interface{} `json:"data"`
}
type eventTenantV1 struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

type eventTenantV2 struct {
	Name string `json:"name"`
	Dob  string `json:"dob"`
}

type errSend struct {
	StatusCode int
}

type testClient struct {
	ctx          context.Context
	cancel       context.CancelFunc
	cacheEventCh <-chan *evcache.Event
	ctl          *client.Ctl
	repo         *client.Repo
	queue        *client.Queue
	recorder     *client.Recorder
	cache        *cache.Client
	httpClient   *http.Client
}

type Result struct {
	Level        string `json:"level"`
	TypeResult   string `json:"Type"`
	EventID      string `json:"Event"`
	Queue        string `json:"Queue"`
	Code         int    `json:"Code"`
	ResultString string `json:"Result"`
	ResultJSON   eventTenantV1
}

func newClient(ctx context.Context, dialer service.GrpcDialer, cliBuilder service.HTTPClient) (*testClient, error) {
	ctl, err := client.NewCtl(ctx, dialer)
	if err != nil {
		return nil, err
	}
	queue, err := client.NewQueue(ctx, dialer)
	if err != nil {
		return nil, err
	}
	repo, err := client.NewRepo(ctx, dialer)
	if err != nil {
		return nil, err
	}
	recorder, err := client.NewRecorder(ctx, dialer)
	if err != nil {
		return nil, err
	}
	cache, err := cache.NewClient(ctx, env.String("cache_listener.addr"), dialer)
	if err != nil {
		return nil, err
	}
	httpClient, err := cliBuilder.NewHTTPClient(env.String("listener.host"))
	if err != nil {
		return nil, err
	}
	ctxCancel, cancel := context.WithCancel(ctx)
	return &testClient{
		ctx:        ctxCancel,
		cancel:     cancel,
		ctl:        ctl,
		queue:      queue,
		repo:       repo,
		recorder:   recorder,
		cache:      cache,
		httpClient: httpClient,
	}, nil
}

type SchemaRefIds struct {
	SchameRef      string
	SchemaRefOk    string
	SchemaRefError string
}

func (s *testClient) close() error {
	errs := make([]error, 0)
	s.cancel()
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
	if err := s.cache.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func NewTestPackage(schemaRefIds SchemaRefIds, runtimeRef string) *pb.JobPackage {
	p := pb.JobPackage{
		ID:     "job_id_1",
		Name:   strptr("job_name_1"),
		Tenant: "tenant_1",
		Queues: []*pb.QueueDef{
			{
				ID:   "queue_id_1",
				Name: strptr("queue_name_1"),
			},
			{
				ID:   "queue_id_1_ok",
				Name: strptr("queue_name_1_ok"),
			},
			{
				ID:   "queue_id_1_error",
				Name: strptr("queue_name_1_error"),
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
					Schema:        &pb.SchemaDef{SchemaRef: schemaRefIds.SchameRef},
				},
				Result: &pb.ResultDef{
					Ok: &pb.EventDef{
						ID:            "event_id_1_ok",
						Name:          strptr("event_name_1_ok"),
						DataType:      pb.DataType_Json,
						SupplierQueue: "queue_id_1_ok",
						Runtime:       "runtime_id_1",
					},
					Error: &pb.EventDef{
						ID:            "event_id_1_error",
						Name:          strptr("event_name_1_error"),
						DataType:      pb.DataType_Json,
						SupplierQueue: "queue_id_1_error",
						Runtime:       "runtime_id_1",
					},
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

func (s *testClient) addPackage(p *pb.JobPackage) error {
	ps, err := s.ctl.Package(s.ctx, p.Tenant, &p.ID)
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

func (s *testClient) updatePackage(p *pb.JobPackage) error {
	ps, err := s.ctl.Package(s.ctx, p.Tenant, &p.ID)
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

func (s *testClient) deletePackage(p *pb.JobPackage) error {
	ps, err := s.ctl.Package(s.ctx, p.Tenant, &p.ID)
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

func (s *testClient) uploadFile(tenant string, fileID string, fileType pb.File_FileType, f io.Reader) error {
	return s.repo.AddFile(s.ctx, tenant, fileID, fileType, f)
}

func (s *testClient) uploadSchemas(p *pb.JobPackage, files map[string][]byte) error {
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

func (s *testClient) uploadRuntimes(p *pb.JobPackage, files map[string][]byte) error {
	for _, e := range p.Runtimes {
		if err := s.uploadFile(p.Tenant,
			e.ModuleRef,
			pb.File_Wasm, bytes.NewReader(files[e.ModuleRef])); err != nil {
			return err
		}
	}
	return nil
}

func (s *testClient) sendEventV1(url *url.URL) (eventTenantV1, error) {
	d := eventTenantV1{"john", "connor", 50}
	ev := event{[]interface{}{d}}
	b, err := json.Marshal(ev)
	if err != nil {
		return eventTenantV1{}, err
	}
	if err := s.sendEvent(url, b); err != nil {
		return eventTenantV1{}, err
	}
	return d, nil
}

func (s *testClient) sendEventV2(url *url.URL) error {
	d := eventTenantV2{"john connor", "01/01/1973"}
	ev := event{[]interface{}{d}}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return s.sendEvent(url, b)
}

func (s *testClient) sendEventMalFormed(url *url.URL) error {
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

func (s *testClient) sendEvent(url *url.URL, e []byte) error {
	r := &http.Request{
		Method: "POST",
		URL:    url,
		Body:   io.NopCloser(bytes.NewReader(e)),
	}
	re, err := s.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer re.Body.Close()
	if re.StatusCode != http.StatusOK {
		return errSend{StatusCode: re.StatusCode}
	}
	return nil
}

func (s *testClient) addTenant(tenant string) error {
	t, err := s.ctl.Tenant(s.ctx, &tenant)
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

func (s *testClient) AllPackages() ([]*pb.JobPackage, error) {
	t, err := s.ctl.AllPackages(s.ctx)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *testClient) dequeue(tenant string, queue string) ([]*pb.QueueItem, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 50*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			res, err := s.queue.Dequeue(s.ctx, tenant, queue)
			if err != nil {
				return nil, err
			}
			if len(res) != 0 {
				return res, nil
			}
			// waiting a bit to avoid bombarding the queue
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (s *testClient) getJobExecutions(pkg *pb.JobPackage, lines int32) ([]Result, error) {
	res, err := s.recorder.JobExecutions(s.ctx, pkg.Tenant, lines)
	if err != nil {
		return nil, err
	}
	result := make([]Result, len(res))
	for idx, r := range res {
		d := json.NewDecoder(strings.NewReader(r))
		result[idx] = Result{}
		if err := d.Decode(&result[idx]); err != nil {
			return nil, err
		}
		if result[idx].TypeResult == strings.ToLower(pb.JobResult_Result.String()) {
			dr := json.NewDecoder(strings.NewReader(result[idx].ResultString))
			m := eventTenantV1{}
			if err := dr.Decode(&m); err != nil {
				return nil, err
			}
			result[idx].ResultJSON = m
		}
	}
	return result, nil
}

func (s *testClient) startRecvCacheEvents() error {
	l, err := s.cache.ListenerForEvents(context.Background())
	if err != nil {
		return err
	}
	s.cacheEventCh = l.C
	return nil
}

func (s *testClient) waitForCacheEvents() error {
	return test.WaitForClosed(context.Background(), s.cacheEventCh)
}
