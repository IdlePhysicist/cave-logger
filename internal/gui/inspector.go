package gui

import (
	"github.com/rivo/tview"
)

type inspector struct {
	*tview.TextView
}

func newInspector(g *Gui) (insp *inspector) {
	insp = &inspector{
		//Frame: tview.NewFrame(tview.NewTextView()),//.SetBorder(true).SetTitle(" Inspector "),
		TextView: tview.NewTextView(),
	}

	insp.SetTitle(` Inspector `).SetTitleAlign(tview.AlignLeft)
	insp.SetBorder(true)
	insp.setInitEntry()
	return
}

func (i *inspector) name() string {
	return `inspector`
}

func (i *inspector) setEntry(text string) {
	i.SetText(text)
}

func (i *inspector) setKeybinding(g *Gui) {}

func (i *inspector) setEntries(g *Gui) {}

func (i *inspector) setInitEntry() {
	i.SetText(`
	     __    __ __      __ 
	|  ||_ |  /  /  \|\/||_  
	|/\||__|__\__\__/|  ||__ `)
}

func (i *inspector) entries(g *Gui) {}

func (i *inspector) updateEntries(g *Gui) {}

func (i *inspector) focus(g *Gui) {
	g.app.SetFocus(i)
}

func (i *inspector) unfocus() {
}

func (i *inspector) setFilterWord(word string) {}