package gui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type stats struct {
	*tview.Table
	active *menu
}

func newStats(g *Gui, menu *menu) (s *stats) {
	s = &stats{
		Table: tview.NewTable().SetSelectable(false, false).SetFixed(3,1),
		active: menu,
	}

	s.SetTitle(` Stats `).SetTitleAlign(tview.AlignLeft)
	s.SetBorder(true)
	s.setEntries(g)
	s.setKeybinding(g)
	return
}

func (s *stats) name() string {
	return `stats`
}

func (s *stats) setEntries(g *Gui) {
	s.entries(g)
	table := s.Clear()

	for i, stat := range g.state.resources.stats {
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

func (s *stats) setKeybinding(g *Gui) {
	s.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		/*switch event.Key() {
		case tcell.KeyEnter:
			g.inspectTrip()
		//case tcell.KeyCtrlR:
		//	t.setEntries(g)
		}*/

		return event
	})
}

func (s *stats) entries(g *Gui) {
	switch g.selectPage(s.active.GetSelection()) {
	case `trips`:
		statSlice := make([]*model.Statistic, 0)
		statSlice = append(
			statSlice,
			&model.Statistic{Name: `Today`, Value: time.Now().Format(`2006-01-02`)},
		)
		g.state.resources.stats = statSlice
	case `caves`:
		g.state.resources.stats, _ = g.db.GetTopCaves()
	case `cavers`:
		g.state.resources.stats, _ = g.db.GetTopCavers()
	}
}

func (s *stats) updateEntries(g *Gui) {}

func (s *stats) focus(g *Gui) {
	s.SetSelectable(true, false)
	g.app.SetFocus(s)
}

func (s *stats) unfocus() {
	s.SetSelectable(false, false)
}

func (s *stats) setFilterWord(word string) {}
