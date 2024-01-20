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

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/protobuf/encoding/prototext"
)

type data struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

type event struct {
	Data []data `json:"data"`
}

type TestData struct {
	Ctx             context.Context
	GrpcDialer      service.GrpcDialer
	TransportSetter service.TranportSetter
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
	client, err := remote.NewControlClient(s.Ctx, s.GrpcDialer)
	if err != nil {
		return err
	}
	ps, err := client.GetPackage(s.Ctx, p.Tenant, &p.ID)
	if err != nil {
		return err
	}
	if len(ps) >= 1 {
		return errors.New("exisys")
	}

	if _, err := client.AddPackage(context.Background(), p); err != nil {
		return err
	}
	return nil
}

func (s *TestData) UpdatePackage(p *pb.JobPackage) error {
	client, err := remote.NewControlClient(s.Ctx, s.GrpcDialer)
	if err != nil {
		return err
	}
	ps, err := client.GetPackage(s.Ctx, p.Tenant, &p.ID)
	if err != nil {
		return err
	}
	if len(ps) == 0 {
		return errors.New("no exisys")
	}

	if err := client.UpdatePackage(context.Background(), p); err != nil {
		return err
	}
	return nil
}

func (s *TestData) DeletePackage(p *pb.JobPackage) error {
	client, err := remote.NewControlClient(s.Ctx, s.GrpcDialer)
	if err != nil {
		return err
	}
	ps, err := client.GetPackage(s.Ctx, p.Tenant, &p.ID)
	if err != nil {
		return err
	}
	if len(ps) == 0 {
		return errors.New("no exisys")
	}

	if err := client.DeletePackage(context.Background(), p); err != nil {
		return err
	}
	return nil
}

func (s *TestData) Uploadfile(tenant string, fileID string, fileType pb.File_FileType, f io.Reader) error {
	client, err := remote.NewRepoClient(s.Ctx, s.GrpcDialer)
	if err != nil {
		return err
	}
	if err = client.AddFile(context.Background(), tenant, fileID, fileType, f); err != nil {
		return err
	}
	return nil
}

func (s *TestData) UploadfileForPackage(p *pb.JobPackage, f io.Reader) error {
	client, err := remote.NewRepoClient(s.Ctx, s.GrpcDialer)
	if err != nil {
		return err
	}
	if err = client.AddFile(context.Background(), p.Tenant, p.Jobs[0].Event.Schema.SchemaRef, pb.File_JsonSchema, f); err != nil {
		return err
	}
	return nil
}

func (s *TestData) SendEvent(url *url.URL) error {
	e := event{[]data{{"sss", "ddd", 1}}}
	j, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if err := s.TransportSetter.Set("fake:1"); err != nil {
		return err
	}
	r := &http.Request{
		Method: "POST",
		URL:    url,
		Body:   io.NopCloser(bytes.NewReader(j)),
	}
	re, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer re.Body.Close()
	if re.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code %d", re.StatusCode)
	}
	return nil
}

func (s *TestData) AddTenant(tenant string) error {
	client, err := remote.NewControlClient(s.Ctx, s.GrpcDialer)
	if err != nil {
		return err
	}
	t, err := client.GetTenant(s.Ctx, &tenant)
	if err != nil {
		return err
	}
	if len(t) >= 1 {
		return errors.New("exist")
	}
	_, err = client.AddTenant(context.Background(), &pb.Tenant{ID: tenant})
	if err != nil {
		return err
	}
	return nil
}

func (s *TestData) GetTenants() ([]*pb.Tenant, error) {
	client, err := remote.NewControlClient(s.Ctx, s.GrpcDialer)
	if err != nil {
		return nil, err
	}
	t, err := client.GetTenants(s.Ctx)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TestData) GetAllPackages() ([]*pb.JobPackage, error) {
	client, err := remote.NewControlClient(s.Ctx, s.GrpcDialer)
	if err != nil {
		return nil, err
	}
	t, err := client.GetAllPackages(s.Ctx)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TestData) GetPackages(tenant string) ([]*pb.JobPackage, error) {
	client, err := remote.NewControlClient(s.Ctx, s.GrpcDialer)
	if err != nil {
		return nil, err
	}
	t, err := client.GetPackages(s.Ctx, tenant)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TestData) Dequeue(tenant string, queue string) ([]*pb.QueueItem, error) {
	i, err := remote.NewQueueClient(context.Background(), s.GrpcDialer)
	if err != nil {
		return nil, err
	}
	res, err := i.Dequeue(context.Background(), tenant, queue)
	if err != nil {
		return nil, err
	}

	return res, nil
}
