package test

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/reflectutil"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/test"
	ctl "github.com/andrescosta/jobico/cmd/ctl/service"
	executor "github.com/andrescosta/jobico/cmd/executor/service"
	listener "github.com/andrescosta/jobico/cmd/listener/service"
	queue "github.com/andrescosta/jobico/cmd/queue/service"
	recorder "github.com/andrescosta/jobico/cmd/recorder/service"
	repo "github.com/andrescosta/jobico/cmd/repo/service"
	pb "github.com/andrescosta/jobico/internal/api/types"
	queuectl "github.com/andrescosta/jobico/internal/queue/controller"
	recorderctl "github.com/andrescosta/jobico/internal/recorder/controller"
	repoctl "github.com/andrescosta/jobico/internal/repo/controller"
)

var (
	//go:embed testdata/schema.json
	schemaV1 []byte

	//go:embed testdata/schema_result_ok.json
	schemaV1Ok []byte

	//go:embed testdata/schema_result_error.json
	schemaV1Error []byte

	//go:embed testdata/echo.wasm
	wasmEcho []byte

	files = map[string][]byte{
		"sch1":       schemaV1,
		"sch1_ok":    schemaV1Ok,
		"sch1_error": schemaV1Error,
		"run1":       wasmEcho,
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
	executor executor.Service
	recorder recorder.Service
}

func (s *JobicoPlatform) ResetQueueService() {
	s.queue = queue.Service{
		Listener: s.conn,
		Dialer:   s.conn,
		Option:   &queuectl.Option{InMemory: true},
	}
}

func NewPlatform() *JobicoPlatform {
	return NewPlatformWithTimeout(*env.Duration("dial.timeout"))
}

func NewPlatformWithTimeout(time time.Duration) *JobicoPlatform {
	conn := service.NewBufConnWithTimeout(time)
	ctl := ctl.Service{
		Listener: conn,
		DBOption: &database.Option{InMemory: true},
		Dialer:   conn,
	}
	queue := queue.Service{
		Listener: conn,
		Dialer:   conn,
		Option:   &queuectl.Option{InMemory: true},
	}
	repo := repo.Service{
		Listener: conn,
		Option:   &repoctl.Option{InMemory: true},
		Dialer:   conn,
	}
	listener := listener.Service{
		Dialer:        conn,
		Listener:      conn,
		ListenerCache: conn,
		ClientBuilder: conn,
	}
	executor := executor.Service{
		Listener:      conn,
		Dialer:        conn,
		ClientBuilder: conn,
	}
	recorder := recorder.Service{
		Listener: conn,
		Option:   &recorderctl.Option{InMemory: true},
		Dialer:   conn,
	}
	return &JobicoPlatform{
		ctl:      ctl,
		conn:     conn,
		queue:    queue,
		repo:     repo,
		listener: listener,
		executor: executor,
		recorder: recorder,
	}
}

type testFn func(*testing.T)

func Test(t *testing.T) {
	tests := []testFn{
		testSunny,
		testEventErrors,
		testQueueDown,
		testErroRepo,
		testErroCtl,
		testErrorInitQueue,
	}
	setEnvVars()
	// defer goleak.VerifyNone(t)
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

func TestStreaming(t *testing.T) {
	t.Skip()
	tests := []testFn{
		testSunnyStreaming,
		testStreamingSchemaUpdate,
		testStreamingDelete,
	}
	setEnvVars()
	//	defer goleak.VerifyNone(t)
	for _, fn := range tests {
		testItFn := fn
		name := reflectutil.FuncName(testItFn)
		name, _ = strings.CutPrefix(name, "Test")
		t.Run(name, func(t *testing.T) {
			testItFn(t)
		})
	}
}

func testSunny(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup(platform.conn)
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.ctl, platform.recorder, platform.repo})
	test.Nil(t, err)
	pkg := addPackageAndFiles(t, cli)
	ps, err := cli.getAllPackages()
	test.Nil(t, err)
	test.NotEmpty(t, ps)
	err = svcGroup.Start([]test.Starter{platform.queue, platform.listener})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.executor})
	test.Nil(t, err)
	_ = sendEventV1AndCheckResultOk(t, pkg, cli)
}

func testSunnyStreaming(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup(platform.conn)
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.ctl, platform.repo, platform.listener, platform.queue, platform.recorder})
	test.Nil(t, err)
	err = cli.startRecvCacheEvents()
	test.Nil(t, err)
	pkg := addPackageAndFiles(t, cli)
	err = cli.waitForCacheEvents()
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.executor})
	test.Nil(t, err)
	_ = sendEventV1AndCheckResultOk(t, pkg, cli)
}

func testStreamingSchemaUpdate(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup(platform.conn)
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.ctl, platform.repo, platform.listener, platform.queue, platform.recorder})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.executor})
	test.Nil(t, err)
	err = cli.startRecvCacheEvents()
	test.Nil(t, err)
	pkg := addPackageAndFiles(t, cli)
	err = cli.waitForCacheEvents()
	test.Nil(t, err)
	url := sendEventV1AndCheckResultOk(t, pkg, cli)
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
	checkExecutionResultOk(t, pkg, cli)
}

func testStreamingDelete(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup(platform.conn)
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.ctl, platform.repo, platform.listener, platform.queue, platform.recorder})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.executor})
	test.Nil(t, err)
	err = cli.startRecvCacheEvents()
	test.Nil(t, err)
	pkg := addPackageAndFiles(t, cli)
	err = cli.waitForCacheEvents()
	test.Nil(t, err)
	url := sendEventV1AndCheckResultOk(t, pkg, cli)
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
	svcGroup := test.NewServiceGroup(platform.conn)
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.ctl, platform.repo})
	test.Nil(t, err)
	pkg := addPackageAndFiles(t, cli)
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.listener, platform.queue})
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup(platform.conn)
	cli, err := newClient(ctx, platform.conn, platform.conn)
	t.Cleanup(func() {
		cancel()
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.ctl, platform.repo, platform.recorder})
	test.Nil(t, err)
	pkg := addPackageAndFiles(t, cli)
	err = svcGroup.Start([]test.Starter{platform.listener, platform.queue})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.executor})
	test.Nil(t, err)
	url := sendEventV1AndCheckResultOk(t, pkg, cli)
	err = svcGroup.Stop(ctx, platform.queue)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.ErrorIs(t, err, errSend{StatusCode: 500})
}

// func TestQueueDownUp(t *testing.T) {
// 	setEnv()
// 	platform := New()
// 	svcGroup := test.NewServiceGroup(platform.conn)
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
	platform := NewPlatformWithTimeout(40 * time.Microsecond)
	svcGroup := test.NewServiceGroup(platform.conn)
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		_ = svcGroup.Stop(context.Background(), platform.listener)
		_ = svcGroup.Stop(context.Background(), platform.queue)
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.repo})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.listener})
	test.NotNil(t, err)
	err = svcGroup.Start([]test.Starter{platform.queue})
	test.NotNil(t, err)
}

func testErrorInitQueue(t *testing.T) {
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup(platform.conn)
	ctx := context.Background()
	cli, err := newClient(ctx, platform.conn, platform.conn)
	t.Cleanup(func() {
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.repo, platform.ctl})
	test.Nil(t, err)
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	svcGroup.Start([]test.Starter{platform.listener})
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.NotNil(t, err)
}

func testErroRepo(t *testing.T) {
	os.Setenv("dial.timeout", (40 * time.Millisecond).String())
	platform := NewPlatform()
	svcGroup := test.NewServiceGroup(platform.conn)
	cli, err := newClient(context.Background(), platform.conn, platform.conn)
	t.Cleanup(func() {
		_ = svcGroup.Stop(context.Background(), platform.listener)
		cleanUp(t, svcGroup, cli)
	})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.ctl})
	test.Nil(t, err)
	_ = addPackage(t, cli)
	err = svcGroup.Start([]test.Starter{platform.queue})
	test.Nil(t, err)
	err = svcGroup.Start([]test.Starter{platform.listener})
	test.NotNil(t, err)
}

func cleanUp(t *testing.T, svcGroup *test.ServiceGroup, cli *testClient) {
	fail := false
	if err := svcGroup.StopAll(); err != nil {
		errw := test.ErrWhileStopping{}
		if errors.As(err, &errw) {
			t.Errorf("error while stopping the service %s: %v", *errw.Starter.Addr(), errw.Err)
		} else {
			t.Errorf("error while stopping %v", err)
		}
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

func addPackageAndFiles(t *testing.T, cli *testClient) *pb.JobPackage {
	pkg := addPackage(t, cli)
	err := cli.uploadSchemas(pkg, files)
	test.Nil(t, err)
	err = cli.uploadRuntimes(pkg, files)
	test.Nil(t, err)
	return pkg
}

func sendEventV1AndCheckResultOk(t *testing.T, pkg *pb.JobPackage, cli *testClient) *url.URL {
	url := sendEventV1(t, pkg, cli)
	checkExecutionResultOk(t, pkg, cli)
	return url
}

func sendEventV1(t *testing.T, pkg *pb.JobPackage, cli *testClient) *url.URL {
	u := fmt.Sprintf("http://listener:1/events/%s/%s", pkg.Tenant, pkg.Jobs[0].Event.ID)
	url, err := url.Parse(u)
	test.Nil(t, err)
	err = cli.sendEventV1(url)
	test.Nil(t, err)
	return url
}

func addPackage(t *testing.T, cli *testClient) *pb.JobPackage {
	pkg := cli.newTestPackage(schemaRefIds{"sch1", "sch1_ok", "sch1_error"}, "run1")
	err := cli.addTenant(pkg.Tenant)
	test.Nil(t, err)
	err = cli.addPackage(pkg)
	test.Nil(t, err)
	return pkg
}

func checkExecutionResultOk(t *testing.T, pkg *pb.JobPackage, cli *testClient) {
	res, err := cli.dequeue(pkg.Tenant, "queue_id_1_ok")
	test.Nil(t, err)
	test.NotEmpty(t, res)
	// l, err := cli.getJobExecutions(pkg, 1)
	// test.Nil(t, err)
	// test.NotEmpty(t, l)
}

func setEnvVars() {
	os.Setenv("dial.timeout", (20 * time.Second).String())
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

	os.Setenv("recorder.addr", "recorder:1")
	os.Setenv("recorder.host", "recorder:1")

	os.Setenv("recorder.host", "recorder:1")
	os.Setenv("recorder.addr", "recorder:1")
	os.Setenv("log.console.enabled", "false")
	os.Setenv("log.file.enabled", "false")
	os.Setenv("executor.timeout", (1 * time.Microsecond).String())
}
