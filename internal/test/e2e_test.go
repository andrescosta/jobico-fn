package test

import (
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

var (
	//go:embed testdata/schema.json
	schemaV1 []byte

	//go:embed testdata/schema_result_ok.json
	schemaV1Ok []byte

	//go:embed testdata/schema_result_error.json
	schemaV1Error []byte

	schemas = map[string][]byte{
		"sch1":       schemaV1,
		"sch1_ok":    schemaV1Ok,
		"sch1_error": schemaV1Error,
	}

	//go:embed testdata/schema_updated.json
	schemaV2  []byte
	schemasV2 = map[string][]byte{
		"sch1_v2": schemaV2,
	}
)

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

type JobicoPlatform struct {
	conn     *service.BufConn
	ctl      ctl.Service
	queue    queue.Service
	repo     repo.Service
	listener listener.Service
}

func (s *JobicoPlatform) ResetQueueService() {
	s.queue = queue.Service{
		Listener: s.conn,
		Dialer:   s.conn,
		Option:   &controller.Option{InMemory: true},
	}
}

func NewPlatform() *JobicoPlatform {
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
	return &JobicoPlatform{
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
	setEnvVars()
	for _, fn := range tests {
		testItFn := fn
		name := reflectutil.FuncName(testItFn)
		name, _ = strings.CutPrefix(name, "Test")
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testItFn(t)
		})
	}
}

func testSunny(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup()
	svcGroup.Start([]test.Starter{platform.ctl, platform.repo})
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	err = cli.addTenant(pkg.Tenant)
	test.Nil(t, err)
	err = cli.addPackage(pkg)
	test.Nil(t, err)
	ps, err := cli.getAllPackages()
	test.Nil(t, err)
	test.NotEmpty(t, ps)
	err = cli.uploadSchemas(pkg, schemas)
	test.Nil(t, err)
	svcGroup.Start([]test.Starter{platform.listener, platform.queue})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.Nil(t, err)
	res, err := cli.dequeue(pkg.Tenant, pkg.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
}

func testSunnyStreaming(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup()
	svcGroup.Start([]test.Starter{platform.ctl, platform.repo, platform.listener, platform.queue})
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = cli.startRecvCacheEvents()
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	err = cli.addTenant(pkg.Tenant)
	test.Nil(t, err)
	err = cli.uploadSchemas(pkg, schemas)
	test.Nil(t, err)
	err = cli.addPackage(pkg)
	test.Nil(t, err)
	err = cli.waitForCacheEvents()
	test.Nil(t, err)
	ps, err := cli.getAllPackages()
	test.Nil(t, err)
	test.NotEmpty(t, ps)
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.Nil(t, err)
	res, err := cli.dequeue(pkg.Tenant, pkg.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
}

func cleanUp(t *testing.T, svcGroup *test.ServiceGroup, cli *client) {
	fail := false
	if err := svcGroup.Stop(); err != nil {
		t.Errorf("error stopping service group %v", err)
		fail = true
	}
	if cli != nil {
		if err := cli.close(); err != nil {
			t.Errorf("error stopping service group %v", err)
			fail = true
		}
	}
	if fail {
		t.FailNow()
	}
}

func testStreamingSchemaUpdate(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup()
	svcGroup.Start([]test.Starter{platform.ctl, platform.repo, platform.listener, platform.queue})
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = cli.startRecvCacheEvents()
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	err = cli.addTenant(pkg.Tenant)
	test.Nil(t, err)
	err = cli.uploadSchemas(pkg, schemas)
	test.Nil(t, err)
	err = cli.addPackage(pkg)
	test.Nil(t, err)
	err = cli.waitForCacheEvents()
	test.Nil(t, err)
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.Nil(t, err)
	res, err := cli.dequeue(pkg.Tenant, pkg.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
	pkg.Jobs[0].Event.Schema.SchemaRef = "sch1_v2"
	err = cli.uploadSchemas(pkg, schemasV2)
	test.Nil(t, err)
	err = cli.updatePackage(pkg)
	test.Nil(t, err)
	err = cli.waitForCacheEvents()
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.NotNil(t, err)
	err = cli.sendEventV2(url)
	test.Nil(t, err)
}

func testStreamingDelete(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup()
	svcGroup.Start([]test.Starter{platform.ctl, platform.repo, platform.listener, platform.queue})
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = cli.startRecvCacheEvents()
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	err = cli.addTenant(pkg.Tenant)
	test.Nil(t, err)
	err = cli.uploadSchemas(pkg, schemas)
	test.Nil(t, err)
	err = cli.addPackage(pkg)
	test.Nil(t, err)
	err = cli.waitForCacheEvents()
	test.Nil(t, err)
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.Nil(t, err)
	res, err := cli.dequeue(pkg.Tenant, pkg.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
	err = cli.deletePackage(pkg)
	test.Nil(t, err)
	err = cli.waitForCacheEvents()
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.NotNil(t, err)
}

func testEventErrors(t *testing.T) {
	setEnvVars()
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup()
	svcGroup.Start([]test.Starter{platform.ctl, platform.repo, platform.listener, platform.queue})
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = cli.startRecvCacheEvents()
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	err = cli.addTenant(pkg.Tenant)
	test.Nil(t, err)
	err = cli.uploadSchemas(pkg, schemas)
	test.Nil(t, err)
	err = cli.addPackage(pkg)
	test.Nil(t, err)
	err = cli.waitForCacheEvents()
	test.Nil(t, err)
	u := fmt.Sprintf("http://listener:1/events/%s/notexist", pkg.Tenant)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.ErrorIs(t, err, errSend{StatusCode: 500})
	u = "http://listener:1/events/fake/notexist"
	url, err = url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.ErrorIs(t, err, errSend{StatusCode: 500})
	u = fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err = url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventMalFormed(url)
	test.ErrorIs(t, err, errSend{StatusCode: 400})
}

func testQueueDown(t *testing.T) {
	setEnvVars()
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup()
	svcGroup.Start([]test.Starter{platform.ctl, platform.repo})
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	err = cli.addTenant(pkg.Tenant)
	test.Nil(t, err)
	err = cli.addPackage(pkg)
	test.Nil(t, err)
	ps, err := cli.getAllPackages()
	test.Nil(t, err)
	test.NotEmpty(t, ps)
	err = cli.uploadSchemas(pkg, schemas)
	test.Nil(t, err)
	svcGroup.Start([]test.Starter{platform.listener, platform.queue})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.Nil(t, err)
	res, err := cli.dequeue(pkg.Tenant, pkg.Queues[0].ID)
	test.Nil(t, err)
	test.NotEmpty(t, res)
	ch, err := svcGroup.StopService(platform.queue)
	test.Nil(t, err)
	err = waitFor(cli.ctx, ch)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.ErrorIs(t, err, errSend{StatusCode: 500})
}

// func TestQueueDownUp(t *testing.T) {
// 	setEnv()
// 	platform := New()
// 	svcGroup := test.NewServiceGroup()
// 	t.Cleanup(func() {
// 		err := svcGroup.Stop()
// 		test.Nil(t, err)
// 	})
// 	svcGroup.AddAndStart([]test.Starter{platform.ctl, platform.repo})
// 	s, err := New(context.Background(), platform.conn, platform.conn)
// 	test.Nil(t, err)
// 	p := s.NewPackage(schemaRefIds{"sch1","",""}, "run1")
// 	err = s.AddTenant(p.Tenant)
// 	test.Nil(t, err)
// 	err = s.AddPackage(p)
// 	test.Nil(t, err)
// 	ps, err := s.GetAllPackages()
// 	test.Nil(t, err)
// 	test.NotEmpty(t, ps)
// 	err = s.UploadfileForPackage(p, schemas)
// 	test.Nil(t, err)
// 	svcGroup.AddAndStart([]test.Starter{platform.listener, platform.queue})
// 	u := fmt.Sprintf("http://listener:1/events/%s/%s", p.Tenant, p.Jobs[0].Event.ID)
// 	url, err := url.Parse(u)
// 	test.Nil(t, err)
// 	err = s.SendEventV1(url)
// 	test.Nil(t, err)
// 	res, err := s.Dequeue(p.Tenant, p.Queues[0].ID)
// 	test.Nil(t, err)
// 	test.NotEmpty(t, res)
// 	err = svcGroup.StopService(platform.queue)
// 	test.Nil(t, err)
// 	_, err = s.GetAllPackages()
// 	test.Nil(t, err)
// 	err = s.SendEventV1(url)
// 	test.ErrorIs(t, err, test.ErrSendEvent{StatusCode: 500})
// 	platform.ResetQueueService()
// 	svcGroup.AddAndStart([]test.Starter{platform.queue})
// 	time.Sleep(5 * time.Second)
// 	err = s.SendEventV1(url)
// 	test.Nil(t, err)
// 	res, err = s.Dequeue(p.Tenant, p.Queues[0].ID)
// 	test.Nil(t, err)
// 	test.NotEmpty(t, res)
// }

func testErroCtl(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup()
	svcGroup.Start([]test.Starter{platform.repo})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Microsecond)
	defer cancel()
	cli, err := newClient(ctx, platform.conn, platform.conn)
	t.Cleanup(func() {
		_ = svcGroup.ResetErrors()
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	svcGroup.StartWithContext(ctx, []test.Starter{platform.listener, platform.queue})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.NotNil(t, err)
	test.NotNil(t, svcGroup.Errors())
}

func testErrorInitQueue(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup()
	svcGroup.Start([]test.Starter{platform.repo, platform.ctl})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Microsecond)
	defer cancel()
	cli, err := newClient(ctx, platform.conn, platform.conn)
	t.Cleanup(func() {
		_ = svcGroup.ResetErrors()
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	svcGroup.StartWithContext(ctx, []test.Starter{platform.listener})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.NotNil(t, err)
	test.NotNil(t, svcGroup.Errors())
}

func testErroRepo(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup()
	svcGroup.Start([]test.Starter{platform.ctl})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Microsecond)
	defer cancel()
	cli, err := newClient(ctx, platform.conn, platform.conn)
	t.Cleanup(func() {
		_ = svcGroup.ResetErrors()
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	svcGroup.StartWithContext(ctx, []test.Starter{platform.listener, platform.queue})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.NotNil(t, err)
	test.NotNil(t, svcGroup.Errors())
}

func setEnvVars() {
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
