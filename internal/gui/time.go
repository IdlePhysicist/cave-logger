package gui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type timeWindow struct {
	*tview.Table
}

func newTimeWindow(g *Gui) (t *timeWindow) {
	t = &timeWindow{
		Table: tview.NewTable().SetSelectable(false, false).SetFixed(3,1),
	}

	t.SetBorder(true)
	t.setEntries(g)
	t.setKeybinding(g)
	return
}

func (t *timeWindow) name() string {
	return `timeWindow`
}

func (t *timeWindow) setEntries(g *Gui) {
	t.entries(g)
	table := t.Clear()

	for i, stat := range g.state.resources.timeWindow {
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

func (t *timeWindow) setKeybinding(g *Gui) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)
		return event
	})
}

func (t *timeWindow) entries(g *Gui) {
	timeSlice := make([]*model.Statistic, 0)
		timeSlice = append(
			timeSlice,
			&model.Statistic{Name: `Today`, Value: time.Now().Format(`2006-01-02`)},
		)
	g.state.resources.timeWindow = timeSlice
}

func (t *timeWindow) updateEntries(g *Gui) {}

func (t *timeWindow) focus(g *Gui) {
	t.SetSelectable(true, false)
	g.app.SetFocus(t)
}

func (t *timeWindow) unfocus() {
	t.SetSelectable(false, false)
}

func (t *timeWindow) setFilterWord(word string) {}
