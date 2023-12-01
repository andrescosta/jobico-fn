package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/andrescosta/goico/pkg/config"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/workflew/api/pkg/remote"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CliApp struct {
	controlClient     *remote.ControlClient
	repoClient        *remote.RepoClient
	recorderClient    *remote.RecorderClient
	infoClients       map[string]*service.GrpcServerInfoClient
	helthCheckClients map[string]*service.GrpcServerHelthCheckClient
	mainView          *tview.Pages
	app               *tview.Application
	lastNode          *tview.TreeNode
	*pb.Environment
	ctxJobResultsGetter    context.Context
	cancelJobResultsGetter context.CancelFunc
}

func main() {
	app := tview.NewApplication().EnableMouse(true)
	err := config.LoadEnvVariables()
	if err != nil {
		os.Exit(-1)
	}
	controlClient, err := remote.NewControlClient()
	if err != nil {
		panic(err)
	}
	repoClient, err := remote.NewRepoClient()
	if err != nil {
		panic(err)
	}
	recorderClient, err := remote.NewRecorderClient()
	if err != nil {
		panic(err)
	}

	cli := &CliApp{
		controlClient:     controlClient,
		repoClient:        repoClient,
		recorderClient:    recorderClient,
		app:               app,
		infoClients:       make(map[string]*service.GrpcServerInfoClient),
		helthCheckClients: make(map[string]*service.GrpcServerHelthCheckClient),
	}
	defer cli.Close()
	cli.render(app)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (c *CliApp) Close() {
	c.controlClient.Close()
	c.repoClient.Close()
	for _, v := range c.infoClients {
		v.Close()
	}
	for _, v := range c.helthCheckClients {
		v.Close()
	}
}

func (c *CliApp) render(app *tview.Application) *tview.Pages {
	pages := tview.NewPages()
	c.mainView = tview.NewPages()
	app.SetRoot(pages, true)

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

	app.SetFocus(menu)

	return pages
}

func (c *CliApp) newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetText(text)
}
func (c *CliApp) renderSideMenu(ctx context.Context) *tview.TreeView {
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
				nl.blur(c)
			}
		}
		ref := n.GetReference().(*node)
		if ref.focus != nil {
			ref.focus(c)
		}
		c.lastNode = n
	})

	return menu
}

func switchToPageIfExists(t *tview.Pages, page string) {
	if t.HasPage(page) {
		t.SwitchToPage(page)
	}
}

func createContentView(content string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true).
		SetWrap(false).
		SetRegions(true)
	fmt.Fprint(textView, content)
	textView.SetBorder(false)
	return textView
}

func switchToEmptyPage(t *tview.Pages) {
	/*if t.HasPage("empty") {
		t.SwitchToPage("empty")
	} else {
		tv := tview.NewTextView()
		fmt.Fprint(tv, "")
		t.AddAndSwitchToPage("empty", tv, true)
	}*/
}
