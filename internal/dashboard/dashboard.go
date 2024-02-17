package dashboard

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andrescosta/goico/pkg/collection"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc"
	"github.com/andrescosta/goico/pkg/service/grpc/svcmeta"
	"github.com/andrescosta/jobico/internal/api/client"
	pb "github.com/andrescosta/jobico/internal/api/types"
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

type Dashboard struct {
	*pb.Environment
	controlCli              *client.Ctl
	repoCli                 *client.Repo
	recorderCli             *client.Recorder
	metadataCli             *client.Metadata
	infoClients             map[string]*svcmeta.InfoClient
	helthCheckClients       map[string]*grpc.HealthCheckClient
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
	dialer                  service.GrpcDialer
}

func New(ctx context.Context, d service.GrpcDialer, name string, sync bool) (*Dashboard, error) {
	loaded, _, err := env.Load(name)
	if err != nil {
		return nil, err
	}
	if !loaded {
		return nil, errors.New(".env files were not loaded")
	}
	controlCli, err := client.NewCtl(ctx, d)
	if err != nil {
		return nil, err
	}
	repoCli, err := client.NewRepo(ctx, d)
	if err != nil {
		return nil, err
	}
	recorderCli, err := client.NewRecorder(ctx, d)
	if err != nil {
		return nil, err
	}
	metadataCli := client.NewMetadata()
	if err != nil {
		return nil, err
	}
	app := tview.NewApplication().EnableMouse(true)
	return &Dashboard{
		controlCli:        controlCli,
		repoCli:           repoCli,
		recorderCli:       recorderCli,
		metadataCli:       metadataCli,
		infoClients:       make(map[string]*svcmeta.InfoClient),
		helthCheckClients: make(map[string]*grpc.HealthCheckClient),
		app:               app,
		sync:              sync,
		dialer:            d,
	}, nil
}

func (c *Dashboard) DebugOn() {
	c.debug = true
}

func (c *Dashboard) Run() error {
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

func (c *Dashboard) Dispose() {
	if err := c.controlCli.Close(); err != nil {
		fmt.Printf("warning: error ctl client: %v\n", err)
	}
	if err := c.repoCli.Close(); err != nil {
		fmt.Printf("warning: error repo client: %v\n", err)
	}
	for _, v := range c.infoClients {
		v.Close()
	}
	for _, v := range c.helthCheckClients {
		if err := v.Close(); err != nil {
			fmt.Printf("warning: error closing health check clients: %v\n", err)
		}
	}
}

func (c *Dashboard) renderUI(ctx context.Context) *tview.Pages {
	// sets the main pages
	pages := tview.NewPages()
	c.mainView = tview.NewPages()
	c.mainView.SetBorderPadding(0, 0, 0, 0)
	c.mainView.SetBorder(true)
	c.mainView.AddPage(emptyPage, buildTextView(""), true, true)
	menu := c.renderSideMenu(ctx)
	c.status = newTextView("")
	c.status.SetTextAlign(tview.AlignCenter)
	helpTxt := "<Esc> - To Exit | <Tab> to switch views | <Arrows> to navigate"
	if c.debug {
		helpTxt = fmt.Sprintf("%s %s", helpTxt, "| <Ctrl-D> for debug info")
	}
	if c.sync {
		helpTxt = fmt.Sprintf("%s %s", helpTxt, "| <Ctrl-P> To stop syncing")
	}
	f := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(c.status, 0, 1, false).
			AddItem(nil, 0, 1, false).
			AddItem(newTextView(helpTxt), 0, 1, false), 0, 1, false)
	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 30).
		SetBorders(true).
		AddItem(newTextView("Jobico Dashboard"), 0, 0, 1, 4, 0, 0, false).
		AddItem(f, 2, 0, 1, 4, 0, 0, false)
	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 1, 0, 1, 1, 0, 0, true).
		AddItem(c.mainView, 1, 1, 1, 3, 0, 0, false)
	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 0, 0, 40, true).
		AddItem(c.mainView, 1, 1, 0, 0, 0, 160, false)
	quitModal := tview.NewModal().SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "Quit" {
				c.app.Stop()
			} else {
				pages.HidePage(quitPageModal)
				c.app.SetFocus(menu)
			}
		})
	c.debugTextView = buildTextView("")
	c.debugTextView.SetWordWrap(true)
	c.debugTextView.SetWrap(true)
	fmt.Fprintf(c.debugTextView, "Debug information for process: %d\n", os.Getppid())
	c.debugTextView.SetBorder(true)
	pages.AddPage(mainPage, grid, true, true)
	pages.AddPage(debugPage, c.debugTextView, true, false)
	// It is important that the last page is always the quit page, so
	// it can appears on top of the others without the need to hide them
	pages.AddPage(quitPageModal, newModal(
		quitModal, 40, 10), true, false)
	c.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		fp, _ := pages.GetFrontPage()
		//exhaustive:ignore
		switch event.Key() {
		case tcell.KeyTAB:
			if fp == mainPage && !quitModal.HasFocus() {
				if menu.HasFocus() {
					c.app.SetFocus(c.mainView)
				} else {
					c.app.SetFocus(menu)
				}
			}
			return nil
		case tcell.KeyEscape, tcell.KeyCtrlC:
			pages.ShowPage(quitPageModal)
			return nil
		case tcell.KeyCtrlD:
			if fp == debugPage {
				pages.SwitchToPage(mainPage)
				c.app.SetFocus(menu)
			} else {
				pages.SwitchToPage(debugPage)
			}
			return nil
		case tcell.KeyCtrlP:
			if c.sync {
				c.showInfo("Streaming stopped", 3*time.Second)
				c.stopStreamingUpdates()
			}
		case tcell.KeyCtrlU:
			if c.debug {
				c.execProtected(func() { panic("testing panic") })
			}
		default:
			return event
		}
		return event
	})
	return pages
}

func (c *Dashboard) renderSideMenu(ctx context.Context) *tview.TreeView {
	e, err := c.controlCli.Environment(ctx)
	if err != nil {
		panic(err)
	}
	ep, err := c.controlCli.AllPackages(ctx)
	if err != nil {
		panic(err)
	}
	fs, err := c.repoCli.AllFilenames(ctx)
	if err != nil {
		panic(err)
	}
	r := renderNode(rootNode(e, ep, fs))
	c.rootTreeNode = r
	menu := tview.NewTreeView()
	menu.SetRoot(r)
	menu.SetCurrentNode(r)
	menu.SetBorder(true)
	menu.SetSelectedFunc(func(n *tview.TreeNode) {
		original := n.GetReference().(*node)
		if len(original.children) > 0 {
			if !original.expanded {
				if n.IsExpanded() {
					c.refreshRootNode(ctx, n)
				}
				iconNodeExpanded := iconPreffixNodesMap[n.IsExpanded()]
				iconNodeClosed := iconPreffixNodesMap[!n.IsExpanded()]
				ns, e := strings.CutPrefix(n.GetText(), iconNodeExpanded)
				if e {
					n.SetText(iconNodeClosed + ns)
					n.SetExpanded(!n.IsExpanded())
				}
			}
		} else if original.selected != nil {
			c.execProtected(func() { original.selected(ctx, c, n) })
		}
	})
	// This function simulates the focus and blur event handlers for the tree's nodes
	menu.SetChangedFunc(func(n *tview.TreeNode) {
		if c.lastNode != nil {
			nl := c.lastNode.GetReference().(*node)
			if nl.blur != nil {
				c.execProtected(func() { nl.blur(ctx, c, c.lastNode, n) })
			}
		}
		ref := n.GetReference().(*node)
		if ref.focus != nil {
			c.execProtected(func() { ref.focus(ctx, c, n) })
		}
		c.lastNode = n
	})
	return menu
}

func (c *Dashboard) refreshRootNode(ctx context.Context, n *tview.TreeNode) {
	original := n.GetReference().(*node)
	switchToEmptyPage(c)
	switch original.rootNodeType {
	case NoRootNode:
		return
	case RootNodePackage:
		ep, err := c.controlCli.AllPackages(ctx)
		if err != nil {
			c.showErrorStr("error refreshing packages data")
			return
		}
		g := packageChildrenNodes(ep)
		original.children = g
		refreshTreeNode(n)
	case RootNodeEnv:
		ep, err := c.controlCli.Environment(ctx)
		if err != nil {
			c.showErrorStr("error refreshing environment data")
			return
		}
		g := environmentChildrenNodes(ep)
		original.children = g
		refreshTreeNode(n)
	case RootNodeFile:
		fs, err := c.repoCli.AllFilenames(ctx)
		if err != nil {
			c.showErrorStr("error refreshing files data")
			return
		}
		g := tenantFileChildrenNodes(fs)
		original.children = g
		refreshTreeNode(n)
	}
}

func (c *Dashboard) startStreamingCtlUpdates(ctx context.Context) error {
	ctx, done := context.WithCancel(ctx)
	c.cancelStreamUpdatesFunc = done
	lp, err := c.controlCli.ListenerForPackageUpdates(ctx)
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
	le, err := c.controlCli.ListenerForEnvironmentUpdates(ctx)
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
	lf, err := c.repoCli.ListenerForRepoUpdates(ctx)
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

func (c *Dashboard) startGettingJobResults(n *tview.TreeNode) {
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
		err := c.recorderCli.StreamJobExecutions(ctxJobResultsGetter, lines, ch)
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

func (c *Dashboard) addNewPackage(p *pb.JobPackage) {
	treeNodePkg, pkgNode := getChidren(RootNodePackage, c.rootTreeNode)
	nn := jobPackageNode(p)
	pkgNode.children = append(pkgNode.children, nn)
	treeNodePkg.AddChild(renderNode(nn))
	if len(pkgNode.children) == 1 {
		reRenderNode(pkgNode, treeNodePkg)
	}
}

func (c *Dashboard) deleteNewPackage(p *pb.JobPackage) {
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

func (c *Dashboard) stopStreamingUpdates() {
	c.cancelStreamUpdatesFunc()
	c.debugInfo("Sync services stopped")
}

func (c *Dashboard) onPanic(e any) {
	txt := fmt.Sprintf("%v", e)
	fmt.Fprintln(c.debugTextView, txt)
	c.showErrorStr("Warning error executing event. Check the debug window.")
}

func (c *Dashboard) execProtected(handler func()) {
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
func (c *Dashboard) showErrorStr(e string, ds ...time.Duration) {
	d := collection.FirstOrDefault(ds, durationError)
	showText(c, c.status, e, tcell.ColorRed, d)
}

func (c *Dashboard) showInfo(e string, ds ...time.Duration) {
	d := collection.FirstOrDefault(ds, durationError)
	showText(c, c.status, e, tcell.ColorBlue, d)
}

func (c *Dashboard) showError(err error, ds ...time.Duration) {
	errStr := collection.UnwrapError(err)[0].Error()
	c.showErrorStr(errStr, collection.FirstOrDefault(ds, durationError))
}

func (c *Dashboard) debugError(err error) {
	log.Err(err)
	fmt.Fprintln(c.debugTextView, err)
}

func (c *Dashboard) debugErrorFromGoRoutine(err error) {
	c.app.QueueUpdateDraw(func() {
		c.debugError(err)
	})
}

func (c *Dashboard) debugInfo(info string) {
	if c.debug {
		fmt.Fprintln(c.debugTextView, info)
	}
}

func (c *Dashboard) debugInfoFromGoRoutine(info string) {
	if c.debug {
		c.app.QueueUpdateDraw(func() {
			c.debugInfo(info)
		})
	}
}
