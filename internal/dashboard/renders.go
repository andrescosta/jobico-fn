package dashboard

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/andrescosta/goico/pkg/service/grpc/svcmeta"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var iconPreffixNodesMap = map[bool]string{
	true:  iconExpanded,
	false: iconContracted,
}

func (c *TApp) renderUI(ctx context.Context) *tview.Pages {
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

func (c *TApp) renderSideMenu(ctx context.Context) *tview.TreeView {
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

func renderNode(target *node) *tview.TreeNode {
	// if target.color == tcell.ColorDefault {
	if len(target.children) > 0 {
		if !target.expanded {
			target.text = renderNodeText(iconContracted, target.text)
		}
		target.color = tcell.ColorGreen
	} else {
		target.color = tcell.ColorWhite
	}
	// }
	node := tview.NewTreeNode(target.text).
		SetExpanded(target.expanded).
		SetReference(target).
		SetColor(target.color)
	for _, child := range target.children {
		node.AddChild(renderNode(child))
	}
	return node
}

func reRenderNode(target *node, tn *tview.TreeNode) {
	// if target.color == tcell.ColorDefault {
	if len(target.children) > 0 {
		if !target.expanded {
			target.text = renderNodeText(iconContracted, target.text)
		}
		target.color = tcell.ColorGreen
	} else {
		newText, f := strings.CutPrefix(target.text, iconContracted)
		if !f {
			newText, _ = strings.CutPrefix(target.text, iconExpanded)
		}
		target.text = newText
		target.color = tcell.ColorWhite
		target.expanded = false
	}
	// }
	tn.SetText(target.text).
		SetExpanded(target.expanded).
		SetReference(target).
		SetColor(target.color)
}

func renderNodeText(icon, text string) string {
	return fmt.Sprintf("%s%s", icon, text)
}

func renderHTTPTableServer(info map[string]string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true)
	table.SetCell(0, 0,
		tview.NewTableCell("Status").
			SetAlign(tview.AlignCenter))
	status := "Unknown"
	table.SetCell(0, 1,
		tview.NewTableCell(status).
			SetAlign(tview.AlignCenter))
	ix := 0
	for k, v := range info {
		table.SetCell(ix+1, 0,
			tview.NewTableCell(k).
				SetAlign(tview.AlignCenter))
		table.SetCell(ix+1, 1,
			tview.NewTableCell(v).
				SetAlign(tview.AlignCenter))
		ix++
	}
	return table
}

func renderGrpcTableServer(infos []*svcmeta.GrpcServerMetadata, s healthpb.HealthCheckResponse_ServingStatus) *tview.Table {
	table := tview.NewTable().
		SetBorders(true)
	table.SetCell(0, 0,
		tview.NewTableCell("Status").
			SetAlign(tview.AlignCenter))
	status := "Not ok"
	if s == healthpb.HealthCheckResponse_SERVING {
		status = "Ok"
	}
	table.SetCell(0, 1,
		tview.NewTableCell(status).
			SetAlign(tview.AlignCenter))
	for ix, info := range infos {
		table.SetCell(ix+1, 0,
			tview.NewTableCell(info.Key).
				SetAlign(tview.AlignCenter))
		table.SetCell(ix+1, 1,
			tview.NewTableCell(info.Value).
				SetAlign(tview.AlignCenter))
	}
	return table
}
