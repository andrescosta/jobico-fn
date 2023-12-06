package tapp

import (
	"fmt"
	"time"

	pb "github.com/andrescosta/jobico/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func buildTextView(text string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true).
		SetWrap(false).
		SetRegions(true)
	fmt.Fprint(textView, text)
	textView.SetBorder(false)
	return textView
}

func switchToPageIfExists(t *tview.Pages, page string) bool {
	if t.HasPage(page) {
		t.SwitchToPage(page)
		return true
	}
	return false
}

func trySwitchToPage(name string, pages *tview.Pages, app *TApp, c func() (tview.Primitive, error)) {
	if !switchToPageIfExists(pages, name) {
		page, err := c()
		if err != nil {
			app.debugError(err)
			app.showError(err)
			return
		}
		if page == nil {
			switchToEmptyPage(app)
			return
		}
		pages.AddAndSwitchToPage(name, page, true)
	}
}

func switchToEmptyPage(app *TApp) {
	app.mainView.SwitchToPage(emptyPage)
}

func showText(app *TApp, status *tview.TextView, text string, color tcell.Color, d time.Duration) {
	status.SetTextColor(color)
	status.SetText(text)
	c := time.NewTimer(d)
	go func() {
		<-c.C
		app.app.QueueUpdateDraw(func() {
			status.SetTextColor(tcell.ColorWhite)
			status.SetText("")
		})
	}()
}

func disableTreeNode(tn *tview.TreeNode) {
	tn.SetColor(tcell.ColorGray)
	n := tn.GetReference().(*node)
	n.selected = func(t *TApp, tn *tview.TreeNode) {}
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

func newTextView(text string) *tview.TextView {
	return tview.NewTextView().
		SetText(text)
}

func getChidren(type1 RootNodeType, tn *tview.TreeNode) (*tview.TreeNode, *node) {
	for _, t := range tn.GetChildren() {
		n := t.GetReference().(*node)
		if n.rootNodeType == type1 {
			return t, n
		}
	}
	return nil, nil
}

func getTenantNode(tenant string, tn *tview.TreeNode) (*tview.TreeNode, *node) {
	for _, t := range tn.GetChildren() {
		n := t.GetReference().(*node)
		e, ok := n.entity.(*pb.TenantFiles)
		if ok && e.TenantId == tenant {
			return t, n
		}
	}
	return nil, nil
}
