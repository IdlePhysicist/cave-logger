package gui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	//"github.com/idlephysicist/cave-logger/internal/model"
)

type statsPeople struct {
	*tview.Table
}

func newStatsPeople(g *Gui) (s *statsPeople) {
	s = &statsPeople{
		Table: tview.NewTable().SetSelectable(false, false).SetFixed(3,1),
	}

	s.SetTitle(` Top Cavers `).SetTitleAlign(tview.AlignLeft)
	s.SetBorder(true)
	s.setEntries(g)
	s.setKeybinding(g)
	return
}

func (s *statsPeople) name() string {
	return `statsPeople`
}

func (s *statsPeople) setEntries(g *Gui) {
	s.entries(g)
	table := s.Clear()

	for i, stat := range g.state.resources.statsPeople {
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

func (s *statsPeople) setKeybinding(g *Gui) {
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

func (s *statsPeople) entries(g *Gui) {
	/*switch g.selectPage(s.active.GetSelection()) {
	case `trips`:
		statSlice := make([]*model.Statistic, 0)
		statSlice = append(
			statSlice,
			&model.Statistic{Name: `Today`, Value: time.Now().Format(`2006-01-02`)},
		)
		g.state.resources.stats = statSlice
	case `caves`:
		g.state.resources.stats, _ = g.db.GetTopCaves()
	case `cavers`:*/
	stats, err := g.db.GetTopCavers()
	if err != nil {
		return
	}
	g.state.resources.statsPeople = stats
}

func (s *statsPeople) updateEntries(g *Gui) {}

func (s *statsPeople) focus(g *Gui) {
	s.SetSelectable(true, false)
	g.app.SetFocus(s)
}

func (s *statsPeople) unfocus() {
	s.SetSelectable(false, false)
}

func (s *statsPeople) setFilterWord(word string) {}
