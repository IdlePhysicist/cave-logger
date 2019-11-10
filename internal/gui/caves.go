package gui

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type caves struct {
	*tview.Table
	caves chan *model.Cave
	filterWord string
}

func newCaves(g *Gui) *caves {
	caves := &caves{
		Table: tview.NewTable().SetSelectable(true, false).Select(0,0).SetFixed(1,1),
	}

	caves.SetTitle(` Caves `).SetTitleAlign(tview.AlignLeft)
	caves.SetBorder(true)
	caves.setEntries(g)
	caves.setKeybinding(g)
	return caves
}

func (c *caves) name() string {
	return `locations`
}

func (c *caves) setKeybinding(g *Gui) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		switch event.Key() {
		case tcell.KeyEnter:
			g.inspectCave()
		case tcell.KeyTAB:
			g.switchPanel(`menu`)
		}

		switch event.Rune() {
		case 'n':
			g.createLocationForm()
		case 'u':
			g.modifyLocationForm()
		case 'd':
			g.deleteLocation()
		}

		return event
	})
}

func (c *caves) entries(g *Gui) {
	caves, err := g.db.GetAllLocations()
	if err != nil {
		return
	}

	g.state.resources.caves = caves	
}

func (c *caves) setEntries(g *Gui) {
	c.entries(g)
	table := c.Clear()

	headers := []string{
		"Name",
		"Region",
		"Country",
		"SRT",
		"Visits",
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

	for i, cave := range g.state.resources.caves {
		table.SetCell(i+1, 0, tview.NewTableCell(cave.Name).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(30).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(cave.Region).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(30).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(cave.Country).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(0).
			SetExpansion(1))

		table.SetCell(i+1, 3, tview.NewTableCell(strconv.FormatBool(cave.SRT)).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(0).
			SetExpansion(1))

		table.SetCell(i+1, 4, tview.NewTableCell(strconv.FormatInt(cave.Visits, 10)).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(0).
			SetExpansion(1))
	}
}

func (c *caves) updateEntries(g *Gui) {}

func (c *caves) focus(g *Gui) {
	c.SetSelectable(true, false)
	g.app.SetFocus(c)
}

func (c *caves) unfocus() {
	c.SetSelectable(false, false)
}

func (c *caves) setFilterWord(word string) {
	c.filterWord = word
}

func (c *caves) monitoringCaves(g *Gui) {
	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
		case <-ticker.C:
			c.updateEntries(g)
		case <-g.state.stopChans["caves"]:
			ticker.Stop()
			break LOOP
		}
	}
}