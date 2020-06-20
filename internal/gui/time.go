package gui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type timeWindow struct {
	*tview.Table
}

func newTimeWindow(g *Gui) (t *timeWindow) {
	t = &timeWindow{
		Table: tview.NewTable().SetSelectable(false, false).SetFixed(3,1),
	}

	t.SetBorder(false)
	t.setEntries()
	//t.setKeybinding(g)
	return
}

func (t *timeWindow) name() string {
	return `timeWindow`
}

func (t *timeWindow) setEntries() {
	table := t.Clear()

	table.SetCell(0, 0, tview.NewTableCell(`Today`).
		SetTextColor(tcell.ColorWhite).
		SetMaxWidth(30).
		SetExpansion(2))

	table.SetCell(0, 1, tview.NewTableCell(time.Now().Format(`2006-01-02`)).
		SetTextColor(tcell.ColorWhite).
		SetMaxWidth(30).
		SetExpansion(1))
}

/*func (t *timeWindow) setKeybinding(g *Gui) {
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
*/
