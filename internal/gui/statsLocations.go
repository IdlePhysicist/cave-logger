package gui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	//"github.com/idlephysicist/cave-logger/internal/model"
)

type statsLocations struct {
	*tview.Table
}

func newStatsLocations(g *Gui) (s *statsLocations) {
	s = &statsLocations{
		Table: tview.NewTable().SetSelectable(false, false).SetFixed(3,1),
	}

	s.SetTitle(` Top Caves `).SetTitleAlign(tview.AlignLeft)
	s.SetBorder(true)
	s.setEntries(g)
	s.setKeybinding(g)
	return
}

func (s *statsLocations) name() string {
	return `statsLocations`
}

func (s *statsLocations) setEntries(g *Gui) {
	s.entries(g)
	table := s.Clear()

	for i, stat := range g.state.resources.statsLocations {
		table.SetCell(i, 0, tview.NewTableCell(stat.Name).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(30).
			SetExpansion(2))

		table.SetCell(i, 1, tview.NewTableCell(stat.Value).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(30).
			SetExpansion(1))
	}
}

func (s *statsLocations) setKeybinding(g *Gui) {
	s.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)
		return event
	})
}

func (s *statsLocations) entries(g *Gui) {
	stats, err := g.db.GetTopCaves()
	if err != nil {
		return
	}
	g.state.resources.statsLocations = stats
}

func (s *statsLocations) updateEntries(g *Gui) {}

func (s *statsLocations) focus(g *Gui) {
	s.SetSelectable(true, false)
	g.app.SetFocus(s)
}

func (s *statsLocations) unfocus() {
	s.SetSelectable(false, false)
}

func (s *statsLocations) setFilterWord(word string) {}
