package gui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	//"github.com/idlephysicist/cave-logger/internal/model"
)

type stats struct {
	*tview.Table
}

func newStats(g *Gui) (s *stats) {
	s = &stats{
		Table: tview.NewTable().SetSelectable(true, false).Select(0,0).SetFixed(1,1),
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

func (s *stats) entries(g *Gui) {}

func (s *stats) updateEntries(g *Gui) {}

func (s *stats) focus(g *Gui) {
	s.SetSelectable(true, false)
	g.app.SetFocus(s)
}

func (s *stats) unfocus() {
	s.SetSelectable(false, false)
}

func (s *stats) setFilterWord(word string) {}