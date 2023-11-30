package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/andrescosta/goico/pkg/config"
	"github.com/andrescosta/workflew/api/pkg/remote"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// const (
// 	lblListenerHost    = "Listener Host"
// 	lblListenerPort    = "Listener Port"
// 	lblQueueHost       = "Queue Host"
// 	lblQueuePort       = "Queue Port"
// 	lblExHost          = "Executor Host"
// 	lblExPort          = "Executor Port"
// 	lblNodeEnvironment = "Environment"
// 	lblNodeDefinitions = "Definitions"
// )

//	type Host struct {
//		hostType string
//		ip       string
//		port     uint32
//	}
//
//	type HostModal struct {
//		visible bool
//		name    string
//		pages   *tview.Pages
//		modal   tview.Primitive
//	}
type CliApp struct {
	controlClient *remote.ControlClient
	repoClient    *remote.RepoClient
	mainView      *tview.Pages
	*pb.Environment
}

func main() {
	app := tview.NewApplication().EnableMouse(true)
	err := config.LoadEnvVariables()
	if err != nil {
		os.Exit(-1)
	}
	cli := &CliApp{
		controlClient: remote.NewControlClient(),
		repoClient:    remote.NewRepoClient(),
	}
	cli.render(app)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (c *CliApp) render(app *tview.Application) *tview.Pages {
	pages := tview.NewPages()
	c.mainView = tview.NewPages()
	app.SetRoot(pages, true)

	menu := c.renderSideMenu(context.Background())

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 30).
		SetBorders(true).
		AddItem(c.newPrimitive("Function as a Processor"), 0, 0, 1, 4, 0, 0, false).
		AddItem(c.newPrimitive("Help"), 2, 0, 1, 4, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 1, 0, 1, 1, 0, 0, false).
		AddItem(c.mainView, 1, 1, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(c.mainView, 1, 1, 1, 3, 0, 120, false)

	pages.AddPage("Mainn", grid, true, true)

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
				target.text = "+ " + target.text
				target.expand = false
				target.color = tcell.ColorGreen
			} else {
				target.color = tcell.ColorWhite
			}
		}
		node := tview.NewTreeNode(target.text).
			SetExpanded(target.expand).
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
	var m = map[bool]string{
		true:  "-",
		false: "+",
	}

	menu.SetSelectedFunc(func(n *tview.TreeNode) {
		original := n.GetReference().(*node)
		if len(original.children) > 0 {
			pref := m[n.IsExpanded()]
			npref := m[!n.IsExpanded()]
			ns, e := strings.CutPrefix(n.GetText(), pref)
			if e {
				n.SetText(npref + ns)
				n.SetExpanded(!n.IsExpanded())
			}
		} else if original.selected != nil {
			original.selected(c)
		}
	})
	menu.SetBlurFunc(func() {
		if c.mainView.HasPage("empty") {
			c.mainView.SwitchToPage("empty")
		} else {
			tv := tview.NewTextView()
			fmt.Fprint(tv, "content")
			c.mainView.AddAndSwitchToPage("empty", tv, true)
		}
	})
	menu.SetFocusFunc(func() {
		/*n := menu.GetCurrentNode()
		if n != nil {
			ref := n.GetReference()
			if ref != nil {
				if reflectico.CanConvert[*node](ref) {
					original := ref.(*node)
					if original.focus != nil {
						original.focus(c)
					}
				}
			}
		}*/
	})

	return menu
}

type node struct {
	text     string
	expand   bool
	selected func(*CliApp)
	focus    func(*CliApp)
	blur     func(*CliApp)
	children []*node
	color    tcell.Color
}

var rootNode = func(e *pb.Environment, j []*pb.JobPackage, r []*pb.TenantFiles) *node {
	return &node{
		text:   "Jobico Manager",
		expand: false,
		children: []*node{
			{text: "Packages", children: jobPackagesNode(j)},
			{text: "Enviroment", children: []*node{
				{text: e.ID, children: []*node{
					{text: "Services", children: servicesNode(e.Services)},
				}},
			}},
			{text: "Files", children: tenantFilesNode(r)},
			{text: "(*) Job Results", color: tcell.ColorGreen, selected: func(c *CliApp) {
				if c.mainView.HasPage("empty") {
					c.mainView.SwitchToPage("empty")
				} else {
					tv := tview.NewTextView()
					fmt.Fprint(tv, "content")
					c.mainView.AddAndSwitchToPage("empty", tv, true)
				}

			}},
		},
	}
}

var serviceNode = func(e *pb.Service) *node {
	return &node{
		text: e.ID,
		children: []*node{
			{text: "Servers", children: serversNode(e.Servers)},
			{text: "Storages", children: storagesNode(e.Storages)},
		},
	}
}

var jobPackageNode = func(e *pb.JobPackage) *node {
	return &node{
		text: e.ID,
	}
}
var serverNode = func(e *pb.Host) *node {
	return &node{
		text: e.Ip + ":" + strconv.Itoa(int(e.Port)),
	}
}

var storageNode = func(s *pb.Storage) *node {
	return &node{
		text: s.ID,
	}
}

var tenantFileNode = func(e *pb.TenantFiles) *node {
	return &node{
		text:     e.TenantId,
		children: filesNode(e.TenantId, e.Files),
	}
}

var fileNode = func(tenant string, e string) *node {
	if strings.HasSuffix(e, ".json") {
		return &node{
			text: e, selected: func(c *CliApp) {
				pageName := tenant + "/" + e
				if c.mainView.HasPage(pageName) {
					c.mainView.SwitchToPage(pageName)
				} else {
					f, err := c.repoClient.GetFile(context.Background(), tenant, e)
					if err == nil {
						pp := RenderJson(string(f))
						c.mainView.AddAndSwitchToPage(pageName, pp, true)
					}
				}

			},
		}
	} else {
		return &node{
			text: e,
		}
	}

}

var serversNode = func(e []*pb.Host) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, serverNode(ee))
	}
	return r
}

var storagesNode = func(e []*pb.Storage) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, storageNode(ee))
	}
	return r
}

var servicesNode = func(e []*pb.Service) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, serviceNode(ee))
	}
	return r
}

var jobPackagesNode = func(e []*pb.JobPackage) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, jobPackageNode(ee))
	}
	return r
}

var filesNode = func(merchant string, e []string) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, fileNode(merchant, ee))
	}
	return r

}

var tenantFilesNode = func(e []*pb.TenantFiles) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, tenantFileNode(ee))
	}
	return r

}

// TextView2 demonstrates the extended text view.
func RenderJson(content string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true).
		SetWrap(false).
		SetRegions(true)
		/*SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				return
			}
			highlights := textView.GetHighlights()
			hasHighlights := len(highlights) > 0
			switch key {
			case tcell.KeyEnter:
				if hasHighlights {
					textView.Highlight()
				} else {
					textView.Highlight("0").
						ScrollToHighlight()
				}
			case tcell.KeyTab:
				if hasHighlights {
					current, _ := strconv.Atoi(highlights[0])
					next := (current + 1) % 9
					textView.Highlight(strconv.Itoa(next)).
						ScrollToHighlight()
				}
			case tcell.KeyBacktab:
				if hasHighlights {
					current, _ := strconv.Atoi(highlights[0])
					next := (current - 1 + 9) % 9
					textView.Highlight(strconv.Itoa(next)).
						ScrollToHighlight()
				}
			}
		})*/
	fmt.Fprint(textView, content)
	textView.SetBorder(true).SetTitle("Content")
	return textView
}

// func (c *CliApp) renderSideMenu() *tview.TreeView {

// 	nodeEnvironment := c.renderEnviromentNode(context.Background())
// 	//nodeJob := c.renderJobDefNode(context.Background())

// 	menu := tview.NewTreeView()
// 	menu.SetRoot(nodeEnvironment).
// 		SetCurrentNode(nodeEnvironment)
// 	return menu
// }

// func (c *CliApp) renderDefinitionsNode(ctx context.Context) *tview.TreeNode {
// 	return tview.NewTreeNode(lblNodeDefinitions).
// 		SetColor(tcell.ColorWhite)
// }

// func (c *CliApp) renderJobDefNode(ctx context.Context) *tview.TreeNode {
// 	/*	nodeEnv := tview.NewTreeNode("+ Job Definitions").
// 			SetColor(tcell.ColorGreen)
// 		nodeEnv.SetSelectable(true)
// 		nodeEnv.SetExpanded(false)

// 		e, err := c.controlClient.GetEnviroment(ctx)
// 		nodeEnvDef := tview.NewTreeNode("+" + e.ID).
// 			SetColor(tcell.ColorWhite).
// 			SetExpanded(false)
// 	*/
// 	return nil
// }

// func (c *CliApp) renderEnviromentNode(ctx context.Context) *tview.TreeNode {
// 	nodeEnv := tview.NewTreeNode("+ Environments").
// 		SetColor(tcell.ColorGreen)
// 	nodeEnv.SetSelectable(true)
// 	nodeEnv.SetExpanded(false)

// 	e, err := c.controlClient.GetEnviroment(ctx)
// 	nodeEnvDef := tview.NewTreeNode("+" + e.ID).
// 		SetColor(tcell.ColorWhite).
// 		SetExpanded(false)

// 	nodeEnv.AddChild(nodeEnvDef)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		if e.Services != nil {
// 			c.Environment = e
// 			nodeSvcs := tview.NewTreeNode("+ Services").
// 				SetColor(tcell.ColorGreen).
// 				SetExpanded(false)
// 			nodeEnvDef.AddChild(nodeSvcs)
// 			for _, s := range e.Services {
// 				c.renderService(nodeSvcs, s)
// 			}
// 		}
// 	}
// 	return nodeEnv
// }

// func (c *CliApp) renderService(n *tview.TreeNode, s *pb.Service) {
// 	tn := tview.NewTreeNode("+" + s.ID).
// 		SetColor(tcell.ColorGreen).
// 		SetExpanded(false)
// 	renderChildren(tn, "Servers", s.Servers)
// 	renderChildren(tn, "Storages", s.Storages)
// 	n.AddChild(tn)
// }

// func renderChildren[T any](parent *tview.TreeNode, name string, hs []T) {
// 	suffix := "+"
// 	if len(hs) == 0 {
// 		suffix = "*"
// 	}
// 	tree := tview.NewTreeNode(fmt.Sprintf("%s %s", suffix, name)).
// 		SetColor(tcell.ColorGreen).
// 		SetExpanded(false)
// 	parent.AddChild(tree)
// 	for _, h := range hs {
// 		id := reflectico.GetFieldString(h, "ID")
// 		child := tview.NewTreeNode(id).
// 			SetColor(tcell.ColorGreen)
// 		tree.AddChild(child)
// 	}
// }
