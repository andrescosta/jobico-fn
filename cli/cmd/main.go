// package main

// import (
// 	"fmt"

// 	"github.com/gdamore/tcell/v2"
// 	"github.com/rivo/tview"
// )

// func main() {
// 	app := tview.NewApplication()
// 	draw(app)
// 	if err := app.Run(); err != nil {
// 		fmt.Printf("Error running application: %s\n", err)
// 	}
// }

// func draw(app *tview.Application) {

// 	//

// 	nodeRoot := tview.NewTreeNode("Components").
// 		SetColor(tcell.ColorGreen)

// 	nodeListener := tview.NewTreeNode("Listeners").
// 		SetColor(tcell.ColorWhite)

// 	nodeRoot.AddChild(nodeListener)
// 	nodeQueues := tview.NewTreeNode("Queues").
// 		SetColor(tcell.ColorGreen)
// 	nodeRoot.AddChild(nodeQueues)
// 	nodeExecutors := tview.NewTreeNode("Executors").
// 		SetColor(tcell.ColorGreen)
// 	nodeRoot.AddChild(nodeExecutors)

// 	tree := tview.NewTreeView().
// 		SetRoot(nodeRoot).
// 		SetCurrentNode(nodeRoot)
// 	tree.SetBorder(true)

// 	//

// 	//servers := tview.NewList().ShowSecondaryText(false)
// 	//servers.SetBorder(true).SetTitle("Servers")

// 	details := tview.NewTable().SetBorders(true)
// 	details.SetBorder(true).SetTitle("Details")
// 	flex := tview.NewFlex().
// 		AddItem(tree, 0, 1, true).
// 		AddItem(details, 0, 1, false)

// 	/*servers.AddItem("10.0.0.1", "", 0, nil)
// 	servers.AddItem("10.0.0.2", "", 0, nil)
// 	servers.AddItem("10.0.0.3", "", 0, nil)
// 	servers.SetChangedFunc(func(i int, tableName string, t string, s rune) {
// 		details.Clear()
// 		adddetails(details, strconv.Itoa(i))
// 	})
// 	servers.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
// 		switch key.Key() {
// 		case tcell.KeyEnter:
// 			return nil
// 		}
// 		return key
// 	})
// 	servers.SetDoneFunc(func() {
// 		app.Stop()
// 	})
// 	*/
// 	pages := tview.NewPages().
// 		AddPage("finderPage", flex, true, true)
// 	app.SetRoot(pages, true)
// 	/*servers.SetCurrentItem(3)
// 	servers.SetCurrentItem(0)*/
// }

// func addcomponents(components *tview.List, s string) {
// 	/*components.AddItem("Queues"+s, "", 0, nil)
// 	components.AddItem("Listener"+s, "", 0, nil)
// 	components.AddItem("Executors"+s, "", 0, nil)*/

// }

// func adddetails(details *tview.Table, d string) {
// 	/*color := tcell.ColorGreenYellow
// 	details.SetCell(0, 0, &tview.TableCell{Text: "Name", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
// 		SetCell(0, 1, &tview.TableCell{Text: "Type", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
// 		SetCell(0, 2, &tview.TableCell{Text: "Size", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
// 		SetCell(0, 3, &tview.TableCell{Text: "Null", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
// 		SetCell(0, 4, &tview.TableCell{Text: "Constraint", Align: tview.AlignCenter, Color: tcell.ColorYellow})

// 	details.SetCell(1, 0, &tview.TableCell{Text: "columnName" + d, Color: color}).
// 		SetCell(1, 1, &tview.TableCell{Text: "dataType" + d, Color: color}).
// 		SetCell(1, 2, &tview.TableCell{Text: "sizeText" + d, Align: tview.AlignRight, Color: color}).
// 		SetCell(1, 3, &tview.TableCell{Text: "isNullable" + d, Align: tview.AlignRight, Color: color}).
// 		SetCell(1, 4, &tview.TableCell{Text: "constraintType.String" + d, Align: tview.AlignLeft, Color: color})
// 	*/
// }

// Demo code for the Grid primitive.
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

func NewHostModal(name string, pages *tview.Pages) *HostModal {
	return &HostModal{
		visible: false,
		name:    name,
		pages:   pages,
	}
}

func (h *HostModal) troggle() {

	if h.visible {
		h.pages.HidePage(h.name)
	} else {
		h.pages.ShowPage(h.name)
	}
	h.visible = !h.visible

}

func main() {
	app := tview.NewApplication().EnableMouse(true)
	addUI(app)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// if hostModal.visible {
// 	i := form.GetFormItemByLabel("Port").(*tview.InputField)
// 	i.SetText(menu.GetCurrentNode().GetText())
// 	//hostModal.modal.
// }

func addUI(app *tview.Application) *tview.Pages {
	pages := tview.NewPages()
	app.SetRoot(pages, true)

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	nodeRoot := tview.NewTreeNode("Components").
		SetColor(tcell.ColorGreen)

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
		SetCurrentNode(nodeRoot)

	main := newPrimitive("Main content")
	//sideBar := newPrimitive("Side Bar")

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Function as a Processor"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Help"), 2, 0, 1, 3, 0, 0, false)

	grid.SetFocusFunc(func() {
		app.SetFocus(menu)
	})

	pages.AddPage("Mainn", grid, true, true)

	hostModal := NewHostModal("host", pages)
	form := formHost(func(f *tview.Form) {
		p := f.GetFormItemByLabel("Port").(*tview.InputField).GetText()
		i := f.GetFormItemByLabel("IP address").(*tview.InputField).GetText()
		addEntityToMenu(menu.GetCurrentNode(), i, p)
		hostModal.troggle()
	}, func() {
		hostModal.troggle()
	})

	hostModal.modal = modal(form, 30, 10)

	pages.AddPage(hostModal.name, hostModal.modal, true, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false)
		//.AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 1, 0, 200, false)
	//AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	menu.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyRune:
			r := key.Rune()
			if r == ' ' {
				hostModal.troggle()
			}
			return nil
		}
		return key
	})
	return pages
}

func addEntityToMenu(n *tview.TreeNode, ip string, port string) {
	tn := tview.NewTreeNode(ip + ":" + port).
		SetColor(tcell.ColorGreen)
	tn.SetReference(&Host{
		hostType: n.GetText(),
		ip:       ip,
		port:     port,
	})
	n.AddChild(tn)
}

func modal(p tview.Primitive, width, height int) tview.Primitive {
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
