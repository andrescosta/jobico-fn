package tui

import (
	"fmt"

	"github.com/rivo/tview"
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

func trySwitchToPage(name string, pages *tview.Pages, c func() (tview.Primitive, error)) {
	if !switchToPageIfExists(pages, name) {
		p, err := c()
		if err != nil {
			// TODO: display errors somehow
		}
		if err == nil {
			pages.AddAndSwitchToPage(name, p, true)
		}
	}
}

func generator[T, Y any](d []T, g func(T) Y) []Y {
	r := make([]Y, 0)
	for _, ee := range d {
		r = append(r, g(ee))
	}
	return r
}

func generatorNamed[T, Y any](name string, d []T, g func(string, T) Y) []Y {
	r := make([]Y, 0)
	for _, ee := range d {
		r = append(r, g(name, ee))
	}
	return r
}
