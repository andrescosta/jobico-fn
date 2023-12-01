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
	c.render()
	if err := c.app.Run(); err != nil {
		return err
	}

	return nil
}

func (c *TApp) Close() {
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
	c.app.SetRoot(pages, true)

	menu := c.renderSideMenu(context.Background())

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(25, 30).
		SetBorders(true).
		AddItem(c.newPrimitive("Function as a Processor"), 0, 0, 1, 4, 0, 0, false).
		AddItem(c.newPrimitive("Help"), 2, 0, 1, 4, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 1, 0, 1, 1, 0, 0, false).
		AddItem(c.mainView, 1, 1, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 0, 0, 40, false).
		AddItem(c.mainView, 1, 1, 0, 0, 0, 160, false)

	pages.AddPage("Mainn", grid, true, true)

	c.app.SetFocus(menu)

	return pages
}

func (c *TApp) newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetText(text)
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
