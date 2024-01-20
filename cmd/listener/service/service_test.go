package service_test

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/reflectutil"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/test"
	ctl "github.com/andrescosta/jobico/cmd/ctl/service"
	listener "github.com/andrescosta/jobico/cmd/listener/service"
	queue "github.com/andrescosta/jobico/cmd/queue/service"
	repo "github.com/andrescosta/jobico/cmd/repo/service"
	"github.com/andrescosta/jobico/cmd/testjobico"
	"github.com/andrescosta/jobico/internal/queue/controller"
	repoctl "github.com/andrescosta/jobico/internal/repo/controller"
)

//go:embed testdata/schema.json
var schema []byte

// Start the server using any port: 127.0.0.1:0?
// Mock queue, ctl, repo
// Solution must be reusable
// Test cases:
// - Sunny
//     - Job, tenant, files exist (no mock) - DONE
// - Errors:
//     - tenant does not exists (no mock)
//     - job does not exists (no mock)
//     - malformed event (no mock)
//     - queue returns an error when queue  (mock with a GRPC service impl or client)
// - Init errors:
//     - cannot connect queue  (no mock)
//     - cannot connect ctl  (no mock)
//     - cannot connect repo  (no mock)
// - Streaming:
//   - sunny
//     - new job package (no mock)  - DONE
//     - update package  (no mock)
//     - delete package  (no mock)
//     - update to json schema (no mock)
//   - multiple listener
//     - no errors (no mock)
//     - unsubscribe (no mock)
//     - communication errors (mock)
//   - connection errors
//     - stopped (mock)
//     - restarted (mock)
//

type Services struct {
	conn     *service.BufConn
	ctl      ctl.Service
	queue    queue.Service
	repo     repo.Service
	listener listener.Service
}

func New() *Services {
	conn := service.NewBufConn()
	ctl := ctl.Service{
		Listener: conn,
		DBOption: &database.Option{InMemory: true},
	}
	queue := queue.Service{
		Listener: conn,
		Dialer:   conn,
		Option:   &controller.Option{InMemory: true},
	}
	repo := repo.Service{
		Listener: conn,
		Option:   &repoctl.Option{InMemory: true},
	}
	listener := listener.Service{
		Dialer:   conn,
		Listener: conn,
	}
	return &Services{
		ctl:      ctl,
		conn:     conn,
		queue:    queue,
		repo:     repo,
		listener: listener,
	}
}

type testFn func(*testing.T)

func Test(t *testing.T) {
	t.Parallel()
	tests := make([]testFn, 0)
	tests = append(tests, sunny)
	tests = append(tests, sunnystreaming)
	for _, test := range tests {
		te := test
		name := reflectutil.FuncName(test)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			te(t)
		})
	}
}

func sunny(t *testing.T) {
	setEnv()
	svcs := New()
	svcGroup := testjobico.NewServiceGroup()
	t.Cleanup(func() {
		if err := svcGroup.Stop(); err != nil {
			t.Errorf("not expected error %v", err)
		}
	})
	svcGroup.AddAndStart([]testjobico.Starter{svcs.ctl, svcs.repo})
	s := testjobico.TestData{
		Ctx:             context.Background(),
		GrpcDialer:      svcs.conn,
		TransportSetter: svcs.conn,
	}
	p := s.NewPackage("sch1", "run1")
	err := s.AddTenant(p.Tenant)
	test.Nil(t, err)
	err = s.AddPackage(p)
	test.Nil(t, err)
	ps, err := s.GetAllPackages()
	test.Nil(t, err)
	test.NotEmpty(t, ps)
	err = s.UploadfileForPackage(p, bytes.NewReader(schema))
	test.Nil(t, err)
	svcGroup.AddAndStart([]testjobico.Starter{svcs.listener, svcs.queue})
	u := fmt.Sprintf("http://fake:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEvent(url)
	test.Nil(t, err)

	res, err := s.Dequeue(p.Tenant, p.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
}

func sunnystreaming(t *testing.T) {
	setEnv()
	svcs := New()
	svcGroup := testjobico.NewServiceGroup()
	t.Cleanup(func() {
		if err := svcGroup.Stop(); err != nil {
			t.Errorf("not expected error %v", err)
		}
	})
	svcGroup.AddAndStart([]testjobico.Starter{svcs.ctl, svcs.repo, svcs.listener, svcs.queue})
	s := testjobico.TestData{
		Ctx:             context.Background(),
		GrpcDialer:      svcs.conn,
		TransportSetter: svcs.conn,
	}
	p := s.NewPackage("sch1", "run1")
	err := s.AddTenant(p.Tenant)
	test.Nil(t, err)
	err = s.UploadfileForPackage(p, bytes.NewReader(schema))
	test.Nil(t, err)
	err = s.AddPackage(p)
	test.Nil(t, err)
	ps, err := s.GetAllPackages()
	test.Nil(t, err)
	test.NotEmpty(t, ps)
	u := fmt.Sprintf("http://fake:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	// We give to listener a bit of time to process the new data
	time.Sleep(30 * time.Millisecond)
	err = s.SendEvent(url)
	test.Nil(t, err)
	res, err := s.Dequeue(p.Tenant, p.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
}

func setEnv() {
	os.Setenv("log.level", "0")
	os.Setenv("log.console.enabled", "true")
	os.Setenv("listener.addr", "fake:1")
	os.Setenv("listener.host", "fake:1")

	os.Setenv("ctl.addr", "fake:2")
	os.Setenv("ctl.host", "fake:2")

	os.Setenv("repo.addr", "fake:3")
	os.Setenv("repo.host", "fake:3")

	os.Setenv("executor.addr", "fake:4")

	os.Setenv("queue.addr", "fake:5")
	os.Setenv("queue.host", "fake:5")

	os.Setenv("recorder.host", "fake:6")
	os.Setenv("recorder.addr", "fake:6")
}
