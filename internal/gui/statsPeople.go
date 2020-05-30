package gui

import (
  "github.com/gdamore/tcell"
  "github.com/rivo/tview"
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
      SetTextColor(tcell.ColorWhite).
      SetMaxWidth(30).
      SetExpansion(2))

    table.SetCell(i, 1, tview.NewTableCell(stat.Value).
      SetTextColor(tcell.ColorWhite).
      SetMaxWidth(30).
      SetExpansion(1))
  }
}

func (s *statsPeople) setKeybinding(g *Gui) {
  s.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    g.setGlobalKeybinding(event)

    return event
  })
}

func (s *statsPeople) entries(g *Gui) {
  stats, err := g.db.GetTopPeople()
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
