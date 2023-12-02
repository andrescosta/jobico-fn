package tapp

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

func createContentView(content string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true).
		SetWrap(false).
		SetRegions(true)
	fmt.Fprint(textView, content)
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
		p, err := c()
		if err != nil {
			log.Err(err)
			errtxt := err.Error()
			e, ok := err.(interface {
				Unwrap() []error
			})
			if ok {
				errtxt = e.Unwrap()[0].Error()
			}
			showText(app.status, errtxt, tcell.ColorRed, 6*time.Second, app)
		}
		if err == nil {
			pages.AddAndSwitchToPage(name, p, true)
		}
	}
}

func showText(status *tview.TextView, text string, color tcell.Color, d time.Duration, app *TApp) {
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
