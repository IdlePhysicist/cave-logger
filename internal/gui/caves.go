package gui

import (
	"strconv"
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type caves struct {
	*cview.Table
	filterCol, filterTerm, filterAction string
}

func newCaves(g *Gui) *caves {
	t := cview.NewTable()
	t.SetScrollBarVisibility(cview.ScrollBarNever)
	t.SetSelectable(true, false)
	t.SetSortClicked(false)
	t.Select(0, 0)
	t.SetFixed(1, 1)

	caves := &caves{Table: t}
	caves.SetBorder(true)
	caves.setEntries(g)

	caves.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		switch event.Key() {
		case tcell.KeyEnter:
			g.state.navigate.update("detail")
			g.inspectCave()
		case tcell.KeyCtrlR:
			caves.setEntries(g)
		}

		switch event.Rune() {
		case 'n':
			g.createLocationForm()
		case 'm':
			g.modifyLocationForm()
		case 'd':
			g.deleteLocation()
		}

		return event
	})
	return caves
}

func (c *caves) name() string {
	return `caves`
}

func (c *caves) entries(g *Gui) {
	caves, err := g.reg.GetAllCaves()
	if err != nil {
		return
	}

	var filteredCaves []*model.Cave
	for _, cave := range caves {
		if c.search(cave) {
			continue
		}
		filteredCaves = append(filteredCaves, cave)
	}
	g.state.resources.locations = filteredCaves
}

func (c *caves) setEntries(g *Gui) {
	c.entries(g)
	c.Clear()

	headers := [][]byte{
		[]byte("Name"),
		[]byte("Region"),
		[]byte("Country"),
		[]byte("SRT"),
		[]byte("Visits"),
		[]byte("Last Visit"),
	}

	for i, header := range headers {
		c.SetCell(0, i, &cview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           cview.AlignLeft,
			Color:           cview.Styles.PrimaryTextColor,
			BackgroundColor: cview.Styles.PrimitiveBackgroundColor,
			Attributes:      tcell.AttrBold,
		})
	}

	var cell *cview.TableCell
	for i, cave := range g.state.resources.locations {
		cell = cview.NewTableCell(cave.Name)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		c.SetCell(i+1, 0, cell)

		cell = cview.NewTableCell(cave.Region)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		c.SetCell(i+1, 1, cell)

		cell = cview.NewTableCell(cave.Country)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(0)
		cell.SetExpansion(1)
		c.SetCell(i+1, 2, cell)

		cell = cview.NewTableCell(yesOrNo(cave.SRT))
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(0)
		cell.SetExpansion(1)
		c.SetCell(i+1, 3, cell)

		cell = cview.NewTableCell(strconv.FormatInt(cave.Visits, 10))
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(0)
		cell.SetExpansion(1)
		c.SetCell(i+1, 4, cell)

		cell = cview.NewTableCell(cave.LastVisit)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(0)
		cell.SetExpansion(1)
		c.SetCell(i+1, 5, cell)
	}

	c.Select(1, 1)
}

func (c *caves) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		c.setEntries(g)
	})
}

func (c *caves) focus(g *Gui) {
	c.SetSelectable(true, false)
	g.app.SetFocus(c)
}

func (c *caves) unfocus() {
	c.SetSelectable(false, false)
}

func (c *caves) setFilter(col, term, action string) {
	c.filterCol = col
	c.filterTerm = term
	c.filterAction = action
}

func (c *caves) monitoringCaves(g *Gui) {
	ticker := time.NewTicker(5 * time.Minute)

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

func (g *Gui) uniqueRegion(input []*model.Cave) []string {
	keys := make(map[string]bool)
	uniq := []string{}

	for _, location := range input {
		if _, value := keys[location.Region]; !value {
			keys[location.Region] = true
			uniq = append(uniq, location.Region)
		}
	}

	return uniq
}

func (g *Gui) uniqueCountry(input []*model.Cave) []string {
	keys := make(map[string]bool)
	uniq := []string{}

	for _, location := range input {
		if _, value := keys[location.Country]; !value {
			keys[location.Country] = true
			uniq = append(uniq, location.Country)
		}
	}

	return uniq
}

func yesOrNo(val bool) string {
	if val {
		return `Y`
	} else {
		return `N`
	}
}

func (c *caves) search(cave *model.Cave) bool {
	// Below *looks* goofy but it all makes sense considering this funciton
	// needs to return false normally!!
	switch c.filterCol {
	case "name":
		if strings.Index(strings.ToLower(cave.Name), c.filterTerm) == -1 {
			return true
		}
		return false
	case "region":
		if strings.Index(strings.ToLower(cave.Region), c.filterTerm) == -1 {
			return true
		}
		return false
	case "country":
		if strings.Index(strings.ToLower(cave.Country), c.filterTerm) == -1 {
			return true
		}
		return false
	case "visits", "":
		switch c.filterAction {
		case "<":
			return cave.Visits > atoi(c.filterTerm)
		case ">":
			return cave.Visits < atoi(c.filterTerm)
		default:
			return false
		}
	default:
		return false
	}
}
