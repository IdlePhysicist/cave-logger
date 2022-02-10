package gui

import "code.rocketnine.space/tslocum/cview"

func newTabBar(g *Gui) (t *cview.TextView) {
	t = cview.NewTextView()
	t.SetDynamicColors(true)
	t.SetRegions(true)
	t.SetWrap(false)
	return
}
