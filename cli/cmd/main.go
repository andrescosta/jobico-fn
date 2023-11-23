package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Host struct {
	hostType string
	ip       string
	port     string
}
type HostModal struct {
	visible bool
	name    string
	pages   *tview.Pages
	modal   tview.Primitive
}

func main() {
	app := tview.NewApplication().EnableMouse(true)
	render(app)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func render(app *tview.Application) *tview.Pages {
	pages := tview.NewPages()
	app.SetRoot(pages, true)

	menu := renderSideMenu()

	main := newPrimitive("Main content")

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 30).
		SetBorders(true).
		AddItem(newPrimitive("Function as a Processor"), 0, 0, 1, 4, 0, 0, false).
		AddItem(newPrimitive("Help"), 2, 0, 1, 4, 0, 0, false)

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

	hostModal := newHostModal("host", pages, menu)

	menu.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyRune:
			r := key.Rune()
			if r == ' ' {
				hostModal.troggleVisibility()
			}
			return nil
		}
		return key
	})
	return pages
}

func renderSideMenu() *tview.TreeView {

	nodeRoot := tview.NewTreeNode("Components").
		SetColor(tcell.ColorGreen)
	nodeRoot.SetSelectable(false)

	nodeListener := tview.NewTreeNode("Listeners").
		SetColor(tcell.ColorWhite)
	nodeRoot.AddChild(nodeListener)

	nodeQueues := tview.NewTreeNode("Queues").
		SetColor(tcell.ColorWhite)
	nodeRoot.AddChild(nodeQueues)

	nodeExecutors := tview.NewTreeNode("Executors").
		SetColor(tcell.ColorWhite)
	nodeRoot.AddChild(nodeExecutors)

	menu := tview.NewTreeView()
	menu.SetRoot(nodeRoot).
		SetCurrentNode(nodeListener)
	return menu
}

func newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetText(text)
}

func newHostModal(name string, pages *tview.Pages, menu *tview.TreeView) *HostModal {
	hostModal := &HostModal{
		visible: false,
		name:    name,
		pages:   pages,
	}

	form := formHost(func(f *tview.Form) {
		p := f.GetFormItemByLabel("Port").(*tview.InputField).GetText()
		i := f.GetFormItemByLabel("IP address").(*tview.InputField).GetText()
		addHostToMenu(menu.GetCurrentNode(), i, p)
		hostModal.troggleVisibility()
	}, func() {
		hostModal.troggleVisibility()
	})

	hostModal.modal = newModal(form, 30, 10)

	pages.AddPage(hostModal.name, hostModal.modal, true, false)

	return hostModal
}

func addHostToMenu(n *tview.TreeNode, ip string, port string) {
	tn := tview.NewTreeNode(ip + ":" + port).
		SetColor(tcell.ColorGreen)
	tn.SetReference(&Host{
		hostType: n.GetText(),
		ip:       ip,
		port:     port,
	})
	n.AddChild(tn)
}

func newModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func formHost(save func(*tview.Form), cancel func()) *tview.Form {
	form := tview.NewForm().
		AddInputField("IP address", "", 15, nil, nil).
		AddInputField("Port", "", 15, nil, nil)

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

func (h *HostModal) troggleVisibility() {
	if h.visible {
		h.pages.HidePage(h.name)
	} else {
		h.pages.ShowPage(h.name)
	}
	h.visible = !h.visible

}
