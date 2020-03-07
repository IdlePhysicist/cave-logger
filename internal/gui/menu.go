package gui

import (
  "github.com/gdamore/tcell"
  "github.com/rivo/tview"
)

type menu struct {
  *tview.Table
}

func newMenu(g *Gui) (m *menu) {
  m = &menu{
    Table: tview.NewTable().SetSelectable(true, false).Select(0,0).SetFixed(1,1),
  }

  m.SetTitle(` Menu `).SetTitleAlign(tview.AlignLeft)
  m.SetBorder(true)
  m.setEntries(g)
  m.setKeybinding(g)
  return
}

func (m *menu) name() string {
  return `menu`
}

func (m *menu) setEntries(g *Gui) {
  m.entries(g)
  table := m.Clear()

  for i, option := range g.state.resources.menu {
    table.SetCell(i, 0, tview.NewTableCell(option).
      SetTextColor(tcell.ColorLightGreen).
      SetMaxWidth(30).
      SetExpansion(1))
  }
}

func (m *menu) setKeybinding(g *Gui) {
  m.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    g.setGlobalKeybinding(event)

    switch event.Key() {
    case tcell.KeyEnter:
      g.goTo(g.selectPage(m.GetSelection()))
    case tcell.KeyTAB:
      g.goTo(g.selectPage(m.GetSelection()))
    }

    return event
  })
}

func (m *menu) entries(g *Gui) {
  options := []string{`Trips`, `People`, `Locations`}
  g.state.resources.menu = options
}

func (m *menu) updateEntries(g *Gui) {}

func (m *menu) focus(g *Gui) {
  m.SetSelectable(true, false)
  g.app.SetFocus(m)
}

func (m *menu) unfocus() {
  m.SetSelectable(false, false)
}

func (m *menu) setFilterWord(word string) {}