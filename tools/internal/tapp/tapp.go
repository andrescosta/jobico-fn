package tapp

import (
	"context"
	"strings"

	"github.com/andrescosta/goico/pkg/config"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/workflew/api/pkg/remote"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TApp struct {
	*pb.Environment
	controlClient          *remote.ControlClient
	repoClient             *remote.RepoClient
	recorderClient         *remote.RecorderClient
	infoClients            map[string]*service.GrpcServerInfoClient
	helthCheckClients      map[string]*service.GrpcServerHelthCheckClient
	app                    *tview.Application
	mainView               *tview.Pages
	lastNode               *tview.TreeNode
	status                 *tview.TextView
	ctxJobResultsGetter    context.Context
	cancelJobResultsGetter context.CancelFunc
}

func New() (*TApp, error) {
	err := config.LoadEnvVariables()
	if err != nil {
		return nil, err
	}
	controlClient, err := remote.NewControlClient()
	if err != nil {
		return nil, err
	}
	repoClient, err := remote.NewRepoClient()
	if err != nil {
		return nil, err
	}
	recorderClient, err := remote.NewRecorderClient()
	if err != nil {
		return nil, err
	}
	app := tview.NewApplication().EnableMouse(true)

	return &TApp{
		controlClient:     controlClient,
		repoClient:        repoClient,
		recorderClient:    recorderClient,
		infoClients:       make(map[string]*service.GrpcServerInfoClient),
		helthCheckClients: make(map[string]*service.GrpcServerHelthCheckClient),
		app:               app,
	}, nil
}

func (c *TApp) Run() error {
	c.app.SetRoot(c.render(), true)
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

func (c *TApp) render() *tview.Pages {
	pages := tview.NewPages()
	c.mainView = tview.NewPages()
	c.mainView.SetBorderPadding(0, 0, 0, 0)
	c.mainView.SetBorder(true)
	menu := c.renderSideMenu(context.Background())
	c.status = newTextView("")
	c.status.SetTextAlign(tview.AlignCenter)

	helpTxt := "<Esc> - To Exit | <Tab> to switch views | <Arrows> to navigate"
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

	const quitPageModal = "quit"
	const mainPage = "main"
	pages.AddPage(mainPage, grid, true, true)
	pages.AddPage(quitPageModal, newModal(
		tview.NewModal().SetText("Do you want to quit the application?").
			AddButtons([]string{"Quit", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Quit" {
					c.app.Stop()
				} else {
					pages.HidePage(quitPageModal)
					c.app.SetFocus(menu)
				}
			}), 40, 10), true, false)

	c.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			if menu.HasFocus() {
				c.app.SetFocus(c.mainView)
			} else {
				c.app.SetFocus(menu)
			}
			return nil
		}
		if event.Key() == tcell.KeyCtrlC ||
			event.Key() == tcell.KeyEscape {
			pages.ShowPage(quitPageModal)
			return nil
		}
		return event
	})

	return pages
}

func (c *TApp) renderSideMenu(ctx context.Context) *tview.TreeView {
	var add func(target *node) *tview.TreeNode
	add = func(target *node) *tview.TreeNode {
		if target.color == tcell.ColorDefault {
			if len(target.children) > 0 {
				if !target.expanded {
					target.text = "+ " + target.text
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
			node.AddChild(add(child))
		}
		return node
	}
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
	r := add(rootNode(e, ep, fs))
	menu := tview.NewTreeView()
	menu.SetRoot(r)
	menu.SetCurrentNode(r)
	var m = map[bool]string{
		true:  "-",
		false: "+",
	}

	menu.SetSelectedFunc(func(n *tview.TreeNode) {
		original := n.GetReference().(*node)
		if len(original.children) > 0 {
			if !original.expanded {
				pref := m[n.IsExpanded()]
				npref := m[!n.IsExpanded()]
				ns, e := strings.CutPrefix(n.GetText(), pref)
				if e {
					n.SetText(npref + ns)
					n.SetExpanded(!n.IsExpanded())
				}
			}
		} else if original.selected != nil {
			original.selected(c, n)
		}
	})
	// This function simulates the focus and blur event handlers for the tree's nodes
	menu.SetChangedFunc(func(n *tview.TreeNode) {
		if c.lastNode != nil {
			nl := c.lastNode.GetReference().(*node)
			if nl.blur != nil {
				nl.blur(c, c.lastNode, n)
			}
		}
		ref := n.GetReference().(*node)
		if ref.focus != nil {
			ref.focus(c, n)
		}
		c.lastNode = n
	})

	return menu
}
