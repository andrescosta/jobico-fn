package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/andrescosta/goico/pkg/config"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/workflew/api/pkg/remote"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	lblListenerHost    = "Listener Host"
	lblListenerPort    = "Listener Port"
	lblQueueHost       = "Queue Host"
	lblQueuePort       = "Queue Port"
	lblExHost          = "Executor Host"
	lblExPort          = "Executor Port"
	lblNodeEnvironment = "Environment"
	lblNodeDefinitions = "Definitions"
)

type Host struct {
	hostType string
	ip       string
	port     uint32
}
type HostModal struct {
	visible bool
	name    string
	pages   *tview.Pages
	modal   tview.Primitive
}
type CliApp struct {
	controlClient *remote.ControlClient
	*pb.Environment
}

func main() {
	app := tview.NewApplication().EnableMouse(true)
	err := config.LoadEnvVariables()
	if err != nil {
		os.Exit(-1)
	}
	cc := remote.NewControlClient(env.GetAsString("ctl.host", ""))
	cli := &CliApp{
		controlClient: cc,
	}
	cli.render(app)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (c *CliApp) render(app *tview.Application) *tview.Pages {
	pages := tview.NewPages()
	app.SetRoot(pages, true)

	menu := c.renderSideMenu()

	main := c.newPrimitive("Main content")

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 30).
		SetBorders(true).
		AddItem(c.newPrimitive("Function as a Processor"), 0, 0, 1, 4, 0, 0, false).
		AddItem(c.newPrimitive("Help"), 2, 0, 1, 4, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 1, 0, 1, 1, 0, 0, false).
		AddItem(main, 1, 1, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 3, 0, 120, false)

	grid.SetFocusFunc(func() {
		app.SetFocus(menu)
	})

	pages.AddPage("Mainn", grid, true, true)

	environmentModal := c.renderEnvironmentModal("host", pages, menu)

	menu.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyRune:
			r := key.Rune()
			if r == ' ' && menu.GetCurrentNode().GetText() == lblNodeEnvironment {
				environmentModal.troggleVisibility()
			}
			return nil
		}
		return key
	})
	return pages
}

func (c *CliApp) renderSideMenu() *tview.TreeView {

	nodeRoot := tview.NewTreeNode("Server").
		SetColor(tcell.ColorGreen)
	nodeRoot.SetSelectable(false)

	nodeEnvironment := c.renderEnviromentNode(context.Background())
	nodeRoot.AddChild(nodeEnvironment)

	nodeMerchants := c.renderDefinitionsNode(context.Background())
	nodeRoot.AddChild(nodeMerchants)

	menu := tview.NewTreeView()
	menu.SetRoot(nodeRoot).
		SetCurrentNode(nodeEnvironment)
	return menu
}

func (c *CliApp) renderDefinitionsNode(ctx context.Context) *tview.TreeNode {
	return tview.NewTreeNode(lblNodeDefinitions).
		SetColor(tcell.ColorWhite)
}

func (c *CliApp) renderEnviromentNode(ctx context.Context) *tview.TreeNode {
	nodeEnvironment := tview.NewTreeNode(lblNodeEnvironment).
		SetColor(tcell.ColorWhite)
	enviro, err := c.controlClient.GetEnviroment(ctx)
	if err != nil {
		// render error
		fmt.Println(err)
	} else {
		if enviro.Servers != nil {
			c.Environment = enviro
		}
		h, ok := enviro.Servers[pb.SrvQueue]
		if ok {
			c.addHostToMenu(nodeEnvironment, h.Ip, h.Port)
		}
		h, ok = enviro.Servers[pb.SrvListener]
		if ok {
			c.addHostToMenu(nodeEnvironment, h.Ip, h.Port)
		}
		h, ok = enviro.Servers[pb.SrvExecutors]
		if ok {
			c.addHostToMenu(nodeEnvironment, h.Ip, h.Port)
		}
	}
	return nodeEnvironment
}

func (c *CliApp) newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetText(text)
}

func (c *CliApp) renderEnvironmentModal(name string, pages *tview.Pages, menu *tview.TreeView) *HostModal {
	hostModal := &HostModal{
		visible: false,
		name:    name,
		pages:   pages,
	}

	form := c.formHost(func(f *tview.Form) {
		var servers map[string]*pb.Host = make(map[string]*pb.Host)
		if c.Environment != nil {
			menu.GetCurrentNode().ClearChildren()
		}
		queueh := &pb.Host{
			Ip:   getFieldString(f, lblQueueHost),
			Port: getFieldInt(f, lblQueuePort),
		}
		c.addHostToMenu(menu.GetCurrentNode(), queueh.Ip, queueh.Port)
		servers[pb.SrvQueue] = queueh

		exh := &pb.Host{
			Ip:   getFieldString(f, lblExHost),
			Port: getFieldInt(f, lblExPort),
		}
		c.addHostToMenu(menu.GetCurrentNode(), exh.Ip, exh.Port)
		servers[pb.SrvExecutors] = exh

		listh := &pb.Host{
			Ip:   getFieldString(f, lblListenerHost),
			Port: getFieldInt(f, lblListenerPort),
		}
		c.addHostToMenu(menu.GetCurrentNode(), listh.Ip, listh.Port)
		servers[pb.SrvListener] = listh

		if c.Environment == nil {
			e, err := c.controlClient.AddEnvironment(context.Background(), &pb.Environment{
				Servers: servers,
			})
			if err != nil {
				fmt.Println(err)
			} else {
				c.Environment = e
			}
		} else {
			c.Environment.Servers = servers
			err := c.controlClient.UpdateEnvironment(context.Background(), c.Environment)
			if err != nil {
				fmt.Println(err)
			}
		}
		hostModal.troggleVisibility()
	}, func() {
		hostModal.troggleVisibility()
	})

	hostModal.modal = c.newModal(form, 40, 17)

	pages.AddPage(hostModal.name, hostModal.modal, true, false)

	return hostModal
}

func getFieldString(f *tview.Form, lbl string) string {
	return f.GetFormItemByLabel(lbl).(*tview.InputField).GetText()
}

func getFieldInt(f *tview.Form, lbl string) uint32 {
	var i uint32
	fmt.Sscan(getFieldString(f, lbl), &i)
	return i
}

func (c *CliApp) addHostToMenu(n *tview.TreeNode, ip string, port uint32) {
	tn := tview.NewTreeNode(ip + ":" + strconv.FormatUint(uint64(port), 10)).
		SetColor(tcell.ColorGreen)
	tn.SetReference(&Host{
		hostType: n.GetText(),
		ip:       ip,
		port:     port,
	})
	n.AddChild(tn)
}

func (c *CliApp) newModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func (c *CliApp) formHost(save func(*tview.Form), cancel func()) *tview.Form {
	form := c.getFormWithFields()
	form.AddButton("Save", func() {
		save(form)
	}).AddButton("Cancel", cancel)

	form.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyEsc:
			cancel()
			return nil
		}
		return key
	})
	form.SetBorder(true)
	return form
}

func (c *CliApp) getFormWithFields() *tview.Form {
	form := tview.NewForm()
	h, ok := c.Environment.Servers[pb.SrvQueue]
	if ok {
		form.AddInputField(lblQueueHost, h.Ip, 15, nil, nil)
		form.AddInputField(lblQueuePort, strconv.FormatUint(uint64(h.Port), 10), 10, nil, nil)
	} else {
		form.AddInputField(lblQueueHost, "", 15, nil, nil)
		form.AddInputField(lblQueuePort, "", 10, nil, nil)

	}
	h, ok = c.Environment.Servers[pb.SrvListener]
	if ok {
		form.AddInputField(lblListenerHost, h.Ip, 15, nil, nil)
		form.AddInputField(lblListenerPort, strconv.FormatUint(uint64(h.Port), 10), 10, nil, nil)
	} else {
		form.AddInputField(lblListenerHost, "", 15, nil, nil)
		form.AddInputField(lblListenerPort, "", 10, nil, nil)

	}
	h, ok = c.Environment.Servers[pb.SrvExecutors]
	if ok {
		form.AddInputField(lblExHost, h.Ip, 15, nil, nil)
		form.AddInputField(lblExPort, strconv.FormatUint(uint64(h.Port), 10), 10, nil, nil)
	} else {
		form.AddInputField(lblExHost, "", 15, nil, nil)
		form.AddInputField(lblExPort, "", 10, nil, nil)
	}
	return form
}

func (h *HostModal) troggleVisibility() {
	if h.visible {
		h.pages.HidePage(h.name)
	} else {
		h.pages.ShowPage(h.name)
	}
	h.visible = !h.visible

}
