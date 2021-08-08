package gui

import (
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	tview "gitlab.com/tslocum/cview"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type caves struct {
	*tview.Table
	caves                 chan *model.Cave
	filterCol, filterTerm string
}

func newCaves(g *Gui) *caves {
	caves := &caves{
		Table: tview.NewTable().
			SetScrollBarVisibility(tview.ScrollBarNever).
			SetSelectable(true, false).
			Select(0, 0).
			SetFixed(1, 1),
	}

	caves.SetTitle(``).SetTitleAlign(tview.AlignLeft)
	caves.SetBorder(true)
	caves.setEntries(g)
	caves.setKeybinding(g)
	return caves
}

func (c *caves) name() string {
	return `caves`
}

func (c *caves) setKeybinding(g *Gui) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		switch event.Key() {
		case tcell.KeyEnter:
			g.state.navigate.update("detail")
			g.inspectCave()
		case tcell.KeyCtrlR:
			c.setEntries(g)
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
	table := c.Clear()

	headers := []string{
		"Name",
		"Region",
		"Country",
		"SRT",
		"Visits",
		"Last Visit",
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tview.Styles.PrimaryTextColor,
			BackgroundColor: tview.Styles.PrimitiveBackgroundColor,
			Attributes:      tcell.AttrBold,
		})
	}

	for i, cave := range g.state.resources.locations {
		table.SetCell(i+1, 0, tview.NewTableCell(cave.Name).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetMaxWidth(30).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(cave.Region).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetMaxWidth(30).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(cave.Country).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetMaxWidth(0).
			SetExpansion(1))

		table.SetCell(i+1, 3, tview.NewTableCell(yesOrNo(cave.SRT)).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetMaxWidth(0).
			SetExpansion(1))

		table.SetCell(i+1, 4, tview.NewTableCell(strconv.FormatInt(cave.Visits, 10)).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetMaxWidth(0).
			SetExpansion(1))

		table.SetCell(i+1, 5, tview.NewTableCell(cave.LastVisit).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetMaxWidth(0).
			SetExpansion(1))
	}
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

func (c *caves) setFilter(col, term string) {
	c.filterCol = col
	c.filterTerm = term
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
	switch c.filterCol {
	case "name", "":
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
	default:
		return false
	}
}
