package tapp

import (
	"context"
	"fmt"
	"time"

	"github.com/andrescosta/goico/pkg/collection"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service/grpc"
	"github.com/andrescosta/goico/pkg/service/grpc/svcmeta"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	durationError  = 6 * time.Second
	emptyPage      = "emptyPage"
	quitPageModal  = "quit"
	mainPage       = "main"
	debugPage      = "debug"
	iconContracted = "+ "
	iconExpanded   = "- "
)

type TApp struct {
	*pb.Environment
	controlClient           *remote.ControlClient
	repoClient              *remote.RepoClient
	recorderClient          *remote.RecorderClient
	metadataClient          *remote.MetadataClient
	infoClients             map[string]*svcmeta.ServerInfoClient
	helthCheckClients       map[string]*grpc.HelthCheckClient
	app                     *tview.Application
	mainView                *tview.Pages
	lastNode                *tview.TreeNode
	rootTreeNode            *tview.TreeNode
	status                  *tview.TextView
	debugTextView           *tview.TextView
	debug                   bool
	cancelJobResultsGetter  context.CancelFunc
	cancelStreamUpdatesFunc context.CancelFunc
	sync                    bool
}

func New(ctx context.Context, name string, sync bool) (*TApp, error) {
	err := env.Load(name)
	if err != nil {
		return nil, err
	}
	controlClient, err := remote.NewControlClient(ctx)
	if err != nil {
		return nil, err
	}
	repoClient, err := remote.NewRepoClient(ctx)
	if err != nil {
		return nil, err
	}
	recorderClient, err := remote.NewRecorderClient(ctx)
	if err != nil {
		return nil, err
	}
	metadataClient := remote.NewMetadataClient()
	if err != nil {
		return nil, err
	}
	app := tview.NewApplication().EnableMouse(true)
	return &TApp{
		controlClient:     controlClient,
		repoClient:        repoClient,
		recorderClient:    recorderClient,
		metadataClient:    metadataClient,
		infoClients:       make(map[string]*svcmeta.ServerInfoClient),
		helthCheckClients: make(map[string]*grpc.HelthCheckClient),
		app:               app,
		sync:              sync,
	}, nil
}
func (c *TApp) DebugOn() {
	c.debug = true
}
func (c *TApp) Run() error {
	ctx, done := context.WithCancel(context.Background())
	defer done()
	c.app.SetRoot(c.renderUI(ctx), true)
	if c.sync {
		if err := c.startStreamingCtlUpdates(ctx); err != nil {
			c.debugErrorFromGoRoutine(err)
		}
	}
	if err := c.app.Run(); err != nil {
		return err
	}
	return nil
}

func (c *TApp) Dispose() {
	c.controlClient.Close()
	c.repoClient.Close()
	for _, v := range c.infoClients {
		v.Close()
	}
	for _, v := range c.helthCheckClients {
		v.Close()
	}
}

func (c *TApp) refreshRootNode(ctx context.Context, n *tview.TreeNode) {
	original := n.GetReference().(*node)
	switchToEmptyPage(c)
	switch original.rootNodeType {
	case NoRootNode:
		return
	case RootNodePackage:
		ep, err := c.controlClient.GetAllPackages(ctx)
		if err != nil {
			c.showErrorStr("error refreshing packages data")
			return
		}
		g := packageChildrenNodes(ep)
		original.children = g
		refreshTreeNode(n)
	case RootNodeEnv:
		ep, err := c.controlClient.GetEnviroment(ctx)
		if err != nil {
			c.showErrorStr("error refreshing environment data")
			return
		}
		g := environmentChildrenNodes(ep)
		original.children = g
		refreshTreeNode(n)
	case RootNodeFile:
		fs, err := c.repoClient.GetAllFileNames(ctx)
		if err != nil {
			c.showErrorStr("error refreshing files data")
			return
		}
		g := tenantFileChildrenNodes(fs)
		original.children = g
		refreshTreeNode(n)
	}
}

func (c *TApp) startStreamingCtlUpdates(ctx context.Context) error {
	ctx, done := context.WithCancel(ctx)
	c.cancelStreamUpdatesFunc = done
	lp, err := c.controlClient.ListenerForPackageUpdates(ctx)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.debugInfoFromGoRoutine("update to package channel stopped")
				return
			case j := <-lp.C:
				c.app.QueueUpdateDraw(func() {
					switch j.Type {
					case pb.UpdateType_New:
						c.addNewPackage(j.Object)
					case pb.UpdateType_Delete:
						c.deleteNewPackage(j.Object)
						switchToEmptyPage(c)
					case pb.UpdateType_Update:
						switchToEmptyPage(c)
						c.deleteNewPackage(j.Object)
						c.addNewPackage(j.Object)
					}
				})
			}
		}
	}()
	le, err := c.controlClient.ListenerForEnvironmentUpdates(ctx)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.debugInfoFromGoRoutine("update to environment channel stopped")
				return
			case e := <-le.C:
				c.app.QueueUpdateDraw(func() {
					p, n := getChidren(RootNodeEnv, c.rootTreeNode)
					ns := environmentChildrenNodes(e.Object)
					n.children = ns
					refreshTreeNode(p)
				})
			}
		}
	}()
	lf, err := c.repoClient.ListenerForRepoUpdates(ctx)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.debugInfoFromGoRoutine("update to file channel stopped")
				return
			case e := <-lf.C:
				c.app.QueueUpdateDraw(func() {
					r, _ := getChidren(RootNodeFile, c.rootTreeNode)
					tr, tn := getTenantNode(e.Object.Tenant, r)
					ns := tenantFileNode(e.Object.Tenant, e.Object.File)
					tn.children = append(tn.children, ns)
					tr.AddChild(renderNode(ns))
				})
			}
		}
	}()
	return nil
}

func (c *TApp) startGettingJobResults(n *tview.TreeNode) {
	var textView *tview.TextView
	lines := int32(5)
	if c.mainView.HasPage("results") {
		c.mainView.SwitchToPage("results")
		_, fp := c.mainView.GetFrontPage()
		textView = fp.(*tview.TextView)
		lines = 0
	} else {
		textView = tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetScrollable(true).
			SetWordWrap(false).
			SetWrap(false).
			SetMaxLines(100)
		c.mainView.AddAndSwitchToPage("results", textView, true)
	}
	ch := make(chan string)
	var ctxJobResultsGetter context.Context
	ctxJobResultsGetter, c.cancelJobResultsGetter = context.WithCancel(context.Background())
	go func(mc <-chan string) {
		for {
			select {
			case <-ctxJobResultsGetter.Done():
				c.debugInfoFromGoRoutine("collector context is done. stopping results collector ")
				return
			case l, ok := <-mc:
				if ok {
					c.app.QueueUpdateDraw(func() {
						fmt.Fprintln(textView, l)
					})
				} else {
					c.debugInfoFromGoRoutine("collector channel is closed. stopping results collector")
					return
				}
			}
		}
	}(ch)
	go func() {
		defer close(ch)
		err := c.recorderClient.StreamJobExecutions(ctxJobResultsGetter, "", lines, ch)
		if err != nil {
			c.debugErrorFromGoRoutine(err)
			c.showErrorStr("Error getting results", 3*time.Second)
			c.app.QueueUpdateDraw(func() {
				onSelectedStopGettingJobResults(ctxJobResultsGetter, c, n)
				disableTreeNode(n)
			})
		}
		c.debugInfoFromGoRoutine("job execution call returned. stopping results collector")
	}()
}

func (c *TApp) addNewPackage(p *pb.JobPackage) {
	treeNodePkg, pkgNode := getChidren(RootNodePackage, c.rootTreeNode)
	nn := jobPackageNode(p)
	pkgNode.children = append(pkgNode.children, nn)
	treeNodePkg.AddChild(renderNode(nn))
	if len(pkgNode.children) == 1 {
		reRenderNode(pkgNode, treeNodePkg)
	}
}

func (c *TApp) deleteNewPackage(p *pb.JobPackage) {
	treeNodePkg, pkgNode := getChidren(RootNodePackage, c.rootTreeNode)
	for _, childNode := range treeNodePkg.GetChildren() {
		refChildNode := (childNode.GetReference().(*node))
		pkg := refChildNode.entity.(*pb.JobPackage)
		if p.ID == pkg.ID {
			treeNodePkg.RemoveChild(childNode)
			pkgNode.removeChild(refChildNode)
			if len(pkgNode.children) == 0 {
				reRenderNode(pkgNode, treeNodePkg)
			}
		}
	}
}
func (c *TApp) stopStreamingUpdates() {
	c.cancelStreamUpdatesFunc()
	c.debugInfo("Sync services stopped")
}

func (c *TApp) onPanic(e any) {
	txt := fmt.Sprintf("%v", e)
	fmt.Fprintln(c.debugTextView, txt)
	c.showErrorStr("Warning error executing event. Check the debug window.")
}

func (c *TApp) execProtected(handler func()) {
	defer func() {
		if p := recover(); p != nil {
			if c.debug {
				c.onPanic(p)
			}
		}
	}()
	handler()
}

// Debug screen updaters
func (c *TApp) showErrorStr(e string, ds ...time.Duration) {
	d := collection.FirstOrDefault(ds, durationError)
	showText(c, c.status, e, tcell.ColorRed, d)
}
func (c *TApp) showError(err error, ds ...time.Duration) {
	errStr := collection.UnwrapError(err)[0].Error()
	c.showErrorStr(errStr, collection.FirstOrDefault(ds, durationError))
}
func (c *TApp) debugError(err error) {
	log.Err(err)
	fmt.Fprintln(c.debugTextView, err)
}
func (c *TApp) debugErrorFromGoRoutine(err error) {
	c.app.QueueUpdateDraw(func() {
		c.debugError(err)
	})
}
func (c *TApp) debugInfo(info string) {
	if c.debug {
		fmt.Fprintln(c.debugTextView, info)
	}
}
func (c *TApp) debugInfoFromGoRoutine(info string) {
	if c.debug {
		c.app.QueueUpdateDraw(func() {
			c.debugInfo(info)
		})
	}
}
