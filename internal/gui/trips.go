package gui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type trips struct {
	*tview.Table
	trips chan *model.Log
	filterWord string
}

func newTrips(g *Gui) *trips {
	trips := &trips{
		Table: tview.NewTable().SetSelectable(true, false).Select(0,0).SetFixed(1,1),
		trips: make(chan *model.Log),
	}

	trips.SetTitle(` Logs `).SetTitleAlign(tview.AlignLeft)
	trips.SetBorder(true)
	trips.setEntries(g)
	trips.setKeybinding(g)
	return trips
}

func (t *trips) name() string {
	return `trips`
}

func (t *trips) setKeybinding(g *Gui) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)


		return event
	})
}

func (t *trips) setEntries(g *Gui) {}

func (t *trips) entries(g *Gui) {}

func (t *trips) updateEntries(g *Gui) {}

func (t *trips) focus(g *Gui) {
	t.SetSelectable(true, false)
	g.app.SetFocus(t)
}

func (t *trips) unfocus() {
	t.SetSelectable(false, false)
}

func (t *trips) setFilterWord(word string) {
	t.filterWord = word
}

func (t *trips) monitoringTrips(g *Gui) {
	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
		case <-ticker.C:
			t.updateEntries(g)
		case <-g.state.stopChans["trips"]:
			ticker.Stop()
			break LOOP
		}
	}
}