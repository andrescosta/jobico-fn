package testjobico

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
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/protobuf/encoding/prototext"
)

type event struct {
	Data []interface{} `json:"data"`
}

type TestData struct {
	Ctx            context.Context
	ctlClient      *remote.ControlClient
	repoClient     *remote.RepoClient
	queueClient    *remote.QueueClient
	recorderClient *remote.RecorderClient
	cacheClient    evcache.CacheServiceClient
	listenerClient *http.Client
}

func New(ctx context.Context, dialer service.GrpcDialer, transport service.HTTPTranporter) (*TestData, error) {
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
	listenerClient, err := NewListenerClient(transport)
	if err != nil {
		return nil, err
	}
	return &TestData{
		Ctx:            ctx,
		ctlClient:      ctl,
		queueClient:    queue,
		repoClient:     repo,
		recorderClient: recorder,
		cacheClient:    cacheClient,
		listenerClient: listenerClient,
	}, nil
}

func NewListenerClient(t service.HTTPTranporter) (*http.Client, error) {
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

func (s *TestData) NewPackageRandom() (*pb.JobPackage, error) {
	k := pb.JobPackage{}
	faker := gofakeit.NewUnlocked(0)
	err := faker.Struct(&k)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func Display(o io.Writer, j []*pb.JobPackage) {
	for _, p := range j {
		_, _ = o.Write([]byte(prototext.Format(p)))
	}
}

func (s *TestData) NewPackage(schemaRef string, runtimeRef string) *pb.JobPackage {
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
		Jobs: []*pb.JobDef{{
			Event: &pb.EventDef{
				ID:            "event_id_1",
				Name:          strptr("event_name_1"),
				DataType:      pb.DataType_Json,
				SupplierQueue: "queue_id_1",
				Runtime:       "runtime_id_1",
				Schema:        &pb.SchemaDef{SchemaRef: schemaRef},
			},
		}},
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

func (s *TestData) AddPackage(p *pb.JobPackage) error {
	ps, err := s.ctlClient.GetPackage(s.Ctx, p.Tenant, &p.ID)
	if err != nil {
		return err
	}
	if len(ps) >= 1 {
		return errors.New("exisys")
	}

	if _, err := s.ctlClient.AddPackage(s.Ctx, p); err != nil {
		return err
	}
	return nil
}

func (s *TestData) UpdatePackage(p *pb.JobPackage) error {
	ps, err := s.ctlClient.GetPackage(s.Ctx, p.Tenant, &p.ID)
	if err != nil {
		return err
	}
	if len(ps) == 0 {
		return errors.New("no exisys")
	}

	if err := s.ctlClient.UpdatePackage(s.Ctx, p); err != nil {
		return err
	}
	return nil
}

func (s *TestData) DeletePackage(p *pb.JobPackage) error {
	ps, err := s.ctlClient.GetPackage(s.Ctx, p.Tenant, &p.ID)
	if err != nil {
		return err
	}
	if len(ps) == 0 {
		return errors.New("no exisys")
	}

	if err := s.ctlClient.DeletePackage(s.Ctx, p); err != nil {
		return err
	}
	return nil
}

func (s *TestData) Uploadfile(tenant string, fileID string, fileType pb.File_FileType, f io.Reader) error {
	return s.repoClient.AddFile(s.Ctx, tenant, fileID, fileType, f)
}

func (s *TestData) UploadfileForPackage(p *pb.JobPackage, f io.Reader) error {
	return s.Uploadfile(p.Tenant, p.Jobs[0].Event.Schema.SchemaRef, pb.File_JsonSchema, f)
}

func (s *TestData) SendEventV1(url *url.URL) error {
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

func (s *TestData) SendEventV2(url *url.URL) error {
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

func (s *TestData) SendEventMalFormed(url *url.URL) error {
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

type ErrSendEvent struct {
	StatusCode int
}

func (e ErrSendEvent) Error() string {
	return fmt.Sprintf("HTTP Status Code:%d", e.StatusCode)
}

func (s *TestData) sendEvent(url *url.URL, e []byte) error {
	r := &http.Request{
		Method: "POST",
		URL:    url,
		Body:   io.NopCloser(bytes.NewReader(e)),
	}
	re, err := s.listenerClient.Do(r)
	if err != nil {
		return err
	}
	defer re.Body.Close()
	if re.StatusCode != http.StatusOK {
		return ErrSendEvent{StatusCode: re.StatusCode}
	}
	return nil
}

func (s *TestData) AddTenant(tenant string) error {
	t, err := s.ctlClient.GetTenant(s.Ctx, &tenant)
	if err != nil {
		return err
	}
	if len(t) >= 1 {
		return errors.New("exist")
	}
	_, err = s.ctlClient.AddTenant(s.Ctx, &pb.Tenant{ID: tenant})
	if err != nil {
		return err
	}
	return nil
}

func (s *TestData) GetTenants() ([]*pb.Tenant, error) {
	t, err := s.ctlClient.GetTenants(s.Ctx)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TestData) GetAllPackages() ([]*pb.JobPackage, error) {
	t, err := s.ctlClient.GetAllPackages(s.Ctx)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TestData) GetPackages(tenant string) ([]*pb.JobPackage, error) {
	t, err := s.ctlClient.GetPackages(s.Ctx, tenant)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TestData) Dequeue(tenant string, queue string) ([]*pb.QueueItem, error) {
	res, err := s.queueClient.Dequeue(s.Ctx, tenant, queue)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *TestData) WaitForCacheUpdates() (chan *evcache.Event, error) {
	r, err := s.cacheClient.Events(s.Ctx, &evcache.Empty{})
	if err != nil {
		return nil, err
	}
	ch := make(chan *evcache.Event)
	ctx, cancel := context.WithTimeout(s.Ctx, 5*time.Second)
	go func() {
		defer func() {
			close(ch)
			cancel()
		}()
		_ = grpchelper.Recv[*evcache.Event](ctx, r, ch)
	}()
	return ch, nil
}
