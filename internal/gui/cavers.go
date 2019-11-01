package gui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type cavers struct {
	*tview.Table
	cavers chan *model.Caver
	filterWord string
}

func newCavers(g *Gui) *cavers {
	cavers := &cavers{
		Table: tview.NewTable().SetSelectable(true, false).Select(0,0).SetFixed(1,1),
	}

	cavers.SetTitle(` Cavers `).SetTitleAlign(tview.AlignLeft)
	cavers.SetBorder(true)
	cavers.setEntries(g)
	cavers.setKeybinding(g)
	return cavers
}

func (c *cavers) name() string {
	return `cavers`
}

func (c *cavers) setKeybinding(g *Gui) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		return event
	})
}

func (c *cavers) setEntries(g *Gui) {}

func (c *cavers) updateEntries(g *Gui) {}

func (c *cavers) entries(g *Gui) {}

func (c *cavers) focus(g *Gui) {
	c.SetSelectable(true, false)
	g.app.SetFocus(c)
}

func (c *cavers) unfocus() {
	c.SetSelectable(false, false)
}

func (c *cavers) setFilterWord(word string) {
	c.filterWord = word
}

func (c *cavers) monitoringCavers(g *Gui) {
	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
		case <-ticker.C:
			c.updateEntries(g)
		case <-g.state.stopChans["cavers"]:
			ticker.Stop()
			break LOOP
		}
	}
}