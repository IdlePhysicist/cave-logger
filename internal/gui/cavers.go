package gui

import (
	"strconv"
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type cavers struct {
	*cview.Table
	filterCol, filterTerm, filterAction string
}

func newCavers(g *Gui) *cavers {
	t := cview.NewTable()
	t.SetScrollBarVisibility(cview.ScrollBarNever)
	t.SetSelectable(true, false)
	t.SetSortClicked(false)
	t.Select(0, 0)
	t.SetFixed(1, 1)

	cavers := &cavers{Table: t}
	cavers.SetBorder(true)
	cavers.setEntries(g)
	cavers.setKeybinding(g)
	return cavers
}

func (c *cavers) name() string {
	return `cavers`
}

func (c *cavers) setKeybinding(g *Gui) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		switch event.Key() {
		case tcell.KeyEnter:
			g.state.navigate.update("detail")
			g.inspectCaver()
		case tcell.KeyCtrlR:
			c.setEntries(g)
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
	c.Clear()

	headers := [][]byte{
		[]byte("Name"),
		[]byte("Club"),
		[]byte("Count"),
		[]byte("Last Trip"),
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
	for i, caver := range g.state.resources.people {
		cell = cview.NewTableCell(caver.Name)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		c.SetCell(i+1, 0, cell)

		cell = cview.NewTableCell(caver.Club)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(0)
		cell.SetExpansion(1)
		c.SetCell(i+1, 1, cell)

		cell = cview.NewTableCell(strconv.FormatInt(caver.Count, 10))
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(0)
		cell.SetExpansion(1)
		c.SetCell(i+1, 2, cell)

		cell = cview.NewTableCell(caver.LastTrip)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(0)
		cell.SetExpansion(1)
		c.SetCell(i+1, 3, cell)
	}

	c.Select(1, 1)
}

func (c *cavers) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		c.setEntries(g)
	})
}

func (c *cavers) entries(g *Gui) {
	cavers, err := g.reg.GetAllCavers()
	if err != nil {
		return
	}

	var filteredCavers []*model.Caver
	for _, caver := range cavers {
		if c.search(caver) {
			continue
		}
		filteredCavers = append(filteredCavers, caver)
	}
	g.state.resources.people = filteredCavers
}

func (c *cavers) focus(g *Gui) {
	c.SetSelectable(true, false)
	g.app.SetFocus(c)
}

func (c *cavers) unfocus() {
	c.SetSelectable(false, false)
}

func (c *cavers) setFilter(col, term, action string) {
	c.filterCol = col
	c.filterTerm = term
	c.filterAction = action
}

func (c *cavers) monitoringCavers(g *Gui) {
	ticker := time.NewTicker(5 * time.Minute)

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

func (c *cavers) search(caver *model.Caver) bool {
	// Below *looks* goofy but it all makes sense considering this funciton
	// needs to return false normally!!
	switch c.filterCol {
	case "name":
		if strings.Index(strings.ToLower(caver.Name), c.filterTerm) == -1 {
			return true
		}
		return false
	case "club":
		if strings.Index(strings.ToLower(caver.Club), c.filterTerm) == -1 {
			return true
		}
		return false
	case "count", "":
		switch c.filterAction {
		case "<":
			return caver.Count > atoi(c.filterTerm)
		case ">":
			return caver.Count < atoi(c.filterTerm)
		default:
			return false
		}
	default:
		return false
	}
}
