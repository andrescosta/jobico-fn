package test

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"strings"
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
	"github.com/andrescosta/jobico/internal/queue/controller"
	repoctl "github.com/andrescosta/jobico/internal/repo/controller"
)

//go:embed testdata/schema.json
var schemaV1 []byte

//go:embed testdata/schema_updated.json
var schemaV2 []byte

// Start the server using any port: 127.0.0.1:0?
// Mock queue, ctl, repo
// Solution must be reusable
// Test cases:
// - Sunny
//     - Job, tenant, files exist (no mock) - DONE
// - Errors:
//     - tenant does not exists (no mock) - DONE
//     - job does not exists (no mock) - DONE
//     - malformed event (no mock) - DONE
//     - queue returns an error when queue  - DONE
// - Init errors:
//     - cannot connect queue  (no mock) - DONE
//     - cannot connect ctl  (no mock) - DONE
//     - cannot connect repo  (no mock) -DONE
// - Streaming:
//   - sunny
//     - new job package (no mock)  - DONE
//     - update package  (no mock)	- DONE
//     - delete package  (no mock)  - DONE
//     - update to json schema (no mock) - DONE
//   - multiple listener
//     - no errors (no mock)
//     - unsubscribe (no mock)
//     - communication errors (mock)
//   - connection errors
//     - stopped (mock) - DONE
//     - restarted (mock) <NOT POSSIBLE>
//

type Services struct {
	conn     *service.BufConn
	ctl      ctl.Service
	queue    queue.Service
	repo     repo.Service
	listener listener.Service
}

func (s *Services) ResetQueueService() {
	s.queue = queue.Service{
		Listener: s.conn,
		Dialer:   s.conn,
		Option:   &controller.Option{InMemory: true},
	}
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
		Dialer:        conn,
		Listener:      conn,
		ListenerCache: conn,
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
	tests := []testFn{
		testSunny,
		testSunnyStreaming,
		testStreamingSchemaUpdate,
		testStreamingDelete,
		testEventErrors,
		testQueueDown,
		testErroRepo,
		testErroCtl,
		testErrorInitQueue,
	}
	setEnv()
	for _, test := range tests {
		testIt := test
		name := reflectutil.FuncName(testIt)
		name, _ = strings.CutPrefix(name, "Test")
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testIt(t)
		})
	}
}

func testSunny(t *testing.T) {
	svcs := New()
	svcGroup := test.NewServiceGroup()
	t.Cleanup(func() {
		err := svcGroup.Stop()
		test.Nil(t, err)
	})
	svcGroup.AddAndStart([]test.Starter{svcs.ctl, svcs.repo})
	s, err := NewClient(context.Background(), svcs.conn, svcs.conn)
	test.Nil(t, err)
	p := s.NewPackage("sch1", "run1")
	err = s.AddTenant(p.Tenant)
	test.Nil(t, err)
	err = s.AddPackage(p)
	test.Nil(t, err)
	ps, err := s.GetAllPackages()
	test.Nil(t, err)
	test.NotEmpty(t, ps)
	err = s.UploadfileForPackage(p, bytes.NewReader(schemaV1))
	test.Nil(t, err)
	svcGroup.AddAndStart([]test.Starter{svcs.listener, svcs.queue})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.Nil(t, err)
	res, err := s.Dequeue(p.Tenant, p.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
}

func testSunnyStreaming(t *testing.T) {
	svcs := New()
	svcGroup := test.NewServiceGroup()
	t.Cleanup(func() {
		err := svcGroup.Stop()
		test.Nil(t, err)
	})
	svcGroup.AddAndStart([]test.Starter{svcs.ctl, svcs.repo, svcs.listener, svcs.queue})
	s, err := NewClient(context.Background(), svcs.conn, svcs.conn)
	ch, _ := s.WaitForCacheUpdates()
	test.Nil(t, err)
	p := s.NewPackage("sch1", "run1")
	err = s.AddTenant(p.Tenant)
	test.Nil(t, err)
	err = s.UploadfileForPackage(p, bytes.NewReader(schemaV1))
	test.Nil(t, err)
	err = s.AddPackage(p)
	test.Nil(t, err)
	<-ch
	test.Nil(t, err)
	ps, err := s.GetAllPackages()
	test.Nil(t, err)
	test.NotEmpty(t, ps)
	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.Nil(t, err)
	res, err := s.Dequeue(p.Tenant, p.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
}

func testStreamingSchemaUpdate(t *testing.T) {
	svcs := New()
	svcGroup := test.NewServiceGroup()
	t.Cleanup(func() {
		err := svcGroup.Stop()
		test.Nil(t, err)
	})
	svcGroup.AddAndStart([]test.Starter{svcs.ctl, svcs.repo, svcs.listener, svcs.queue})
	s, err := NewClient(context.Background(), svcs.conn, svcs.conn)
	ch, _ := s.WaitForCacheUpdates()
	test.Nil(t, err)
	p := s.NewPackage("sch1_v1", "run1")
	err = s.AddTenant(p.Tenant)
	test.Nil(t, err)
	err = s.UploadfileForPackage(p, bytes.NewReader(schemaV1))
	test.Nil(t, err)
	err = s.AddPackage(p)
	test.Nil(t, err)
	<-ch
	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.Nil(t, err)
	res, err := s.Dequeue(p.Tenant, p.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
	p.Jobs[0].Event.Schema.SchemaRef = "sch1_v2"
	err = s.UploadfileForPackage(p, bytes.NewReader(schemaV2))
	test.Nil(t, err)
	err = s.UpdatePackage(p)
	test.Nil(t, err)
	<-ch
	err = s.SendEventV1(url)
	test.NotNil(t, err)
	err = s.SendEventV2(url)
	test.Nil(t, err)
}

func testStreamingDelete(t *testing.T) {
	svcs := New()
	svcGroup := test.NewServiceGroup()
	t.Cleanup(func() {
		err := svcGroup.Stop()
		test.Nil(t, err)
	})
	svcGroup.AddAndStart([]test.Starter{svcs.ctl, svcs.repo, svcs.listener, svcs.queue})
	s, err := NewClient(context.Background(), svcs.conn, svcs.conn)
	ch, _ := s.WaitForCacheUpdates()
	test.Nil(t, err)
	p := s.NewPackage("sch1_v1", "run1")
	err = s.AddTenant(p.Tenant)
	test.Nil(t, err)
	err = s.UploadfileForPackage(p, bytes.NewReader(schemaV1))
	test.Nil(t, err)
	err = s.AddPackage(p)
	test.Nil(t, err)
	<-ch
	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.Nil(t, err)
	res, err := s.Dequeue(p.Tenant, p.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
	err = s.DeletePackage(p)
	test.Nil(t, err)
	<-ch
	err = s.SendEventV1(url)
	test.NotNil(t, err)
}

func testEventErrors(t *testing.T) {
	setEnv()
	svcs := New()
	svcGroup := test.NewServiceGroup()
	t.Cleanup(func() {
		err := svcGroup.Stop()
		test.Nil(t, err)
	})
	svcGroup.AddAndStart([]test.Starter{svcs.ctl, svcs.repo, svcs.listener, svcs.queue})
	s, err := NewClient(context.Background(), svcs.conn, svcs.conn)
	ch, _ := s.WaitForCacheUpdates()
	test.Nil(t, err)
	p := s.NewPackage("sch1", "run1")
	err = s.AddTenant(p.Tenant)
	test.Nil(t, err)
	err = s.UploadfileForPackage(p, bytes.NewReader(schemaV1))
	test.Nil(t, err)
	err = s.AddPackage(p)
	test.Nil(t, err)
	<-ch
	u := fmt.Sprintf("http://listener:1/events/%s/notexist", p.Tenant)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.ErrorIs(t, err, ErrSendEvent{StatusCode: 500})
	u = "http://listener:1/events/fake/notexist"
	url, err = url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.ErrorIs(t, err, ErrSendEvent{StatusCode: 500})
	u = fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err = url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventMalFormed(url)
	test.ErrorIs(t, err, ErrSendEvent{StatusCode: 400})
}

func testQueueDown(t *testing.T) {
	setEnv()
	svcs := New()
	svcGroup := test.NewServiceGroup()
	t.Cleanup(func() {
		err := svcGroup.Stop()
		test.Nil(t, err)
	})
	svcGroup.AddAndStart([]test.Starter{svcs.ctl, svcs.repo})
	s, err := NewClient(context.Background(), svcs.conn, svcs.conn)
	test.Nil(t, err)
	p := s.NewPackage("sch1", "run1")
	err = s.AddTenant(p.Tenant)
	test.Nil(t, err)
	err = s.AddPackage(p)
	test.Nil(t, err)
	ps, err := s.GetAllPackages()
	test.Nil(t, err)
	test.NotEmpty(t, ps)
	err = s.UploadfileForPackage(p, bytes.NewReader(schemaV1))
	test.Nil(t, err)
	svcGroup.AddAndStart([]test.Starter{svcs.listener, svcs.queue})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.Nil(t, err)
	res, err := s.Dequeue(p.Tenant, p.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
	ch, err := svcGroup.StopService(svcs.queue)
	test.Nil(t, err)
	<-ch
	err = s.SendEventV1(url)
	test.ErrorIs(t, err, ErrSendEvent{StatusCode: 500})
}

// func TestQueueDownUp(t *testing.T) {
// 	setEnv()
// 	svcs := New()
// 	svcGroup := test.NewServiceGroup()
// 	t.Cleanup(func() {
// 		err := svcGroup.Stop()
// 		test.Nil(t, err)
// 	})
// 	svcGroup.AddAndStart([]test.Starter{svcs.ctl, svcs.repo})
// 	s, err := New(context.Background(), svcs.conn, svcs.conn)
// 	test.Nil(t, err)
// 	p := s.NewPackage("sch1", "run1")
// 	err = s.AddTenant(p.Tenant)
// 	test.Nil(t, err)
// 	err = s.AddPackage(p)
// 	test.Nil(t, err)
// 	ps, err := s.GetAllPackages()
// 	test.Nil(t, err)
// 	test.NotEmpty(t, ps)
// 	err = s.UploadfileForPackage(p, bytes.NewReader(schemaV1))
// 	test.Nil(t, err)
// 	svcGroup.AddAndStart([]test.Starter{svcs.listener, svcs.queue})
// 	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
// 	url, err := url.Parse(u)
// 	test.Nil(t, err)
// 	err = s.SendEventV1(url)
// 	test.Nil(t, err)
// 	res, err := s.Dequeue(p.Tenant, p.Queues[0].ID)
// 	test.Nil(t, err)
// 	test.NotEmpty(t, res)
// 	err = svcGroup.StopService(svcs.queue)
// 	test.Nil(t, err)
// 	_, err = s.GetAllPackages()
// 	test.Nil(t, err)
// 	err = s.SendEventV1(url)
// 	test.ErrorIs(t, err, test.ErrSendEvent{StatusCode: 500})
// 	svcs.ResetQueueService()
// 	svcGroup.AddAndStart([]test.Starter{svcs.queue})
// 	time.Sleep(5 * time.Second)
// 	err = s.SendEventV1(url)
// 	test.Nil(t, err)
// 	res, err = s.Dequeue(p.Tenant, p.Queues[0].ID)
// 	test.Nil(t, err)
// 	test.NotEmpty(t, res)
// }

func testErroCtl(t *testing.T) {
	svcs := New()
	svcGroup := test.NewServiceGroup()
	t.Cleanup(func() {
		err := svcGroup.Stop()
		test.Nil(t, err)
	})
	svcGroup.AddAndStart([]test.Starter{svcs.repo})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Microsecond)
	defer cancel()
	s, err := NewClient(ctx, svcs.conn, svcs.conn)
	test.Nil(t, err)
	p := s.NewPackage("sch1", "run1")
	svcGroup.AddAndStartWithContext(ctx, []test.Starter{svcs.listener, svcs.queue})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.NotNil(t, err)
	test.NotNil(t, svcGroup.Errors())
	_ = svcGroup.ResetErrors()
}

func testErrorInitQueue(t *testing.T) {
	svcs := New()
	svcGroup := test.NewServiceGroup()
	t.Cleanup(func() {
		err := svcGroup.Stop()
		test.Nil(t, err)
	})
	svcGroup.AddAndStart([]test.Starter{svcs.repo, svcs.ctl})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Microsecond)
	defer cancel()
	s, err := NewClient(ctx, svcs.conn, svcs.conn)
	test.Nil(t, err)
	p := s.NewPackage("sch1", "run1")
	svcGroup.AddAndStartWithContext(ctx, []test.Starter{svcs.listener})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.NotNil(t, err)
	test.NotNil(t, svcGroup.Errors())
	_ = svcGroup.ResetErrors()
}

func testErroRepo(t *testing.T) {
	svcs := New()
	svcGroup := test.NewServiceGroup()
	t.Cleanup(func() {
		err := svcGroup.Stop()
		test.Nil(t, err)
	})
	svcGroup.AddAndStart([]test.Starter{svcs.ctl})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Microsecond)
	defer cancel()
	s, err := NewClient(ctx, svcs.conn, svcs.conn)
	test.Nil(t, err)
	p := s.NewPackage("sch1", "run1")
	svcGroup.AddAndStartWithContext(ctx, []test.Starter{svcs.listener, svcs.queue})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = s.SendEventV1(url)
	test.NotNil(t, err)
	//	test.NotNil(t, svcGroup.Errors())
	_ = svcGroup.ResetErrors()
}

func setEnv() {
	os.Setenv("log.level", "0")
	os.Setenv("log.console.enabled", "true")
	os.Setenv("listener.addr", "listener:1")
	os.Setenv("listener.host", "listener:1")
	os.Setenv("cache_listener.addr", "cache_listener:1")

	os.Setenv("ctl.addr", "ctl:1")
	os.Setenv("ctl.host", "ctl:1")

	os.Setenv("repo.addr", "repo:1")
	os.Setenv("repo.host", "repo:1")

	os.Setenv("executor.addr", "exec:1")

	os.Setenv("queue.addr", "queue:1")
	os.Setenv("queue.host", "queue:1")

	os.Setenv("recorder.host", "recorder:1")
	os.Setenv("recorder.addr", "recorder:1")
}
