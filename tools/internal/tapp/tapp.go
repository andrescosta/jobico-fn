package tapp

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andrescosta/goico/pkg/collection"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service/grpc"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	durationError = 6 * time.Second

	emptyPage = "emptyPage"

	quitPageModal = "quit"

	mainPage = "main"

	debugPage = "debug"

	iconContracted = "+"

	iconExpanded = "-"
)

type TApp struct {
	*pb.Environment

	controlClient *remote.ControlClient

	repoClient *remote.RepoClient

	recorderClient *remote.RecorderClient

	metadataClient *remote.MetadataClient

	infoClients map[string]*grpc.ServerInfoClient

	helthCheckClients map[string]*grpc.HelthCheckClient

	app *tview.Application

	mainView *tview.Pages

	lastNode *tview.TreeNode

	root *tview.TreeNode

	status *tview.TextView

	debugTextView *tview.TextView

	debug bool

	cancelJobResultsGetter context.CancelFunc

	cancelStreamUpdatesFunc context.CancelFunc

	sync bool
}

func New(ctx context.Context, sync bool) (*TApp, error) {
	err := env.Populate()

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

		controlClient: controlClient,

		repoClient: repoClient,

		recorderClient: recorderClient,

		metadataClient: metadataClient,

		infoClients: make(map[string]*grpc.ServerInfoClient),

		helthCheckClients: make(map[string]*grpc.HelthCheckClient),

		app: app,

		sync: sync,
	}, nil
}

func (c *TApp) DebugOn() {
	c.debug = true
}

func (c *TApp) Run() error {
	ctx, done := context.WithCancel(context.Background())

	defer done()

	c.app.SetRoot(c.render(ctx), true)

	if c.sync {
		if err := c.startStreamingUpdates(ctx); err != nil {
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

func (c *TApp) render(ctx context.Context) *tview.Pages {
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
		helpTxt = fmt.Sprintf("%s %s", helpTxt, "| <Ctrl-D> for debug info | <Ctrl-P> To stop streaming")
	}

	f := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(c.status, 0, 1, false).
			AddItem(nil, 0, 1, false).
			AddItem(newTextView(helpTxt), 0, 1, false), 0, 1, false)

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(25, 30).
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
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
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
			if c.debug {
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

func renderNode(target *node) *tview.TreeNode {
	if target.color == tcell.ColorDefault {
		if len(target.children) > 0 {
			if !target.expanded {
				target.text = fmt.Sprintf("%s %s", iconContracted, target.text)
			}

			target.color = tcell.ColorGreen
		} else {
			target.color = tcell.ColorWhite
		}
	}

	node := tview.NewTreeNode(target.text).
		SetExpanded(target.expanded).
		SetReference(target).
		SetColor(target.color)

	for _, child := range target.children {
		node.AddChild(renderNode(child))
	}

	return node
}

func (c *TApp) renderSideMenu(ctx context.Context) *tview.TreeView {
	e, err := c.controlClient.GetEnviroment(ctx)

	if err != nil {
		panic(err)
	}

	ep, err := c.controlClient.GetAllPackages(ctx)

	if err != nil {
		panic(err)
	}

	fs, err := c.repoClient.GetAllFileNames(ctx)

	if err != nil {
		panic(err)
	}

	r := renderNode(rootNode(e, ep, fs))

	c.root = r

	menu := tview.NewTreeView()

	menu.SetRoot(r)

	menu.SetCurrentNode(r)

	menu.SetBorder(true)

	var m = map[bool]string{

		true: iconExpanded,

		false: iconContracted,
	}

	menu.SetSelectedFunc(func(n *tview.TreeNode) {
		original := n.GetReference().(*node)

		if len(original.children) > 0 {
			if !original.expanded {
				if n.IsExpanded() {
					c.refreshRootNode(ctx, n)
				}

				pref := m[n.IsExpanded()]

				npref := m[!n.IsExpanded()]

				ns, e := strings.CutPrefix(n.GetText(), pref)

				if e {
					n.SetText(npref + ns)

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

func (c *TApp) startStreamingUpdates(ctx context.Context) error {
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
					p, n := getChidren(RootNodeEnv, c.root)

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
					r, _ := getChidren(RootNodeFile, c.root)

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

func (c *TApp) addNewPackage(p *pb.JobPackage) {
	r, n := getChidren(RootNodePackage, c.root)

	nn := jobPackageNode(p)

	n.children = append(n.children, nn)

	r.AddChild(renderNode(nn))
}

func (c *TApp) deleteNewPackage(p *pb.JobPackage) {
	r, np := getChidren(RootNodePackage, c.root)

	for _, ns := range r.GetChildren() {
		n := (ns.GetReference().(*node))

		t := n.entity.(*pb.JobPackage)

		if p.ID == t.ID {
			r.RemoveChild(ns)

			np.removeChild(n)
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

func refreshTreeNode(n *tview.TreeNode) {
	n.ClearChildren()

	for _, child := range n.GetReference().(*node).children {
		n.AddChild(renderNode(child))
	}
}

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
