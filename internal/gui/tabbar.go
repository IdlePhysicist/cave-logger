package gui

import (
	"github.com/rivo/tview"
)

func newTabBar(g *Gui) *tview.TextView {
	return tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)/*.
		SetHighlightedFunc(func(added, removed, remaining []string) {
			g.pages.SwitchToPage(added[0])
		})*/
}
