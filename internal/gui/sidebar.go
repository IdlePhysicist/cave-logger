package gui

import (
	"github.com/rivo/tview"

	//"github.com/idlephysicist/cave-logger/internal/model"
)

type sidebar struct {
	*tview.Table
}

func newSidebar(g *Gui) (bar *sidebar) {
	bar = &sidebar{
		Table: tview.NewTable().SetSelectable(true, false).Select(0,0).SetFixed(1,1),
	}

	bar.SetTitle(` Menu `).SetTitleAlign(tview.AlignLeft)
	bar.SetBorder(true)
	bar.setEntries(g)
	bar.setKeybinding(g)
	return
}

func (b *sidebar) name() string {
	return `sidebar`
}

func (b *sidebar) setEntries(g *Gui) {}

func (b *sidebar) setKeybinding(g *Gui) {}

func (b *sidebar) entries(g *Gui) {}

func (b *sidebar) updateEntries(g *Gui) {}

func (b *sidebar) focus(g *Gui) {
	b.SetSelectable(true, false)
	g.app.SetFocus(b)
}

func (b *sidebar) unfocus() {
	b.SetSelectable(false, false)
}

func (b *sidebar) setFilterWord(word string) {}