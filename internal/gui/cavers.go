package gui

import (
  "strconv"
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
  return `people`
}

func (c *cavers) setKeybinding(g *Gui) {
  c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    g.setGlobalKeybinding(event)

    switch event.Key() {
    case tcell.KeyEnter:
      g.inspectPerson()
    case tcell.KeyTAB:
      g.switchPanel(`menu`)
    }

    switch event.Rune() {
    case 'n':
      g.createPersonForm()
    case 'm':
      g.modifyPersonForm()
    case 'd':
      g.deletePerson()
    }

    return event
  })
}

func (c *cavers) setEntries(g *Gui) {
  c.entries(g)
  table := c.Clear()

  headers := []string{
    "Name",
    "Club",
    "Count",
  }

  for i, header := range headers {
    table.SetCell(0, i, &tview.TableCell{
      Text:            header,
      NotSelectable:   true,
      Align:           tview.AlignLeft,
      Color:           tcell.ColorWhite,
      BackgroundColor: tcell.ColorDefault,
      Attributes:      tcell.AttrBold,
    })
  }

  for i, caver := range g.state.resources.people {
    table.SetCell(i+1, 0, tview.NewTableCell(caver.Name).
      SetTextColor(tcell.ColorLightGreen).
      SetMaxWidth(30).
      SetExpansion(1))

    table.SetCell(i+1, 1, tview.NewTableCell(caver.Club).
      SetTextColor(tcell.ColorLightGreen).
      SetMaxWidth(0).
      SetExpansion(1))

    table.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatInt(caver.Count, 10)).
      SetTextColor(tcell.ColorLightGreen).
      SetMaxWidth(0).
      SetExpansion(1))
  }
}

func (c *cavers) updateEntries(g *Gui) {
  g.app.QueueUpdateDraw(func() {
    c.setEntries(g)
  })
}

func (c *cavers) entries(g *Gui) {
  cavers, err := g.db.GetAllPeople()
  if err != nil {
    return
  }

  g.state.resources.people = cavers	
}

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

func (g *Gui) uniqueClubs(input []*model.Caver) []string {
  keys := make(map[string]bool)
  uniq := []string{}

  for _, person := range input {
    if _, value := keys[person.Club]; !value {
      keys[person.Club] = true
      uniq = append(uniq, person.Club)
    }
  }

  return uniq
}