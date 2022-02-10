package gui

import (
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type trips struct {
	*cview.Table
	filterCol, filterTerm, filterAction string
}

func newTrips(g *Gui) *trips {
	t := cview.NewTable()
	t.SetScrollBarVisibility(cview.ScrollBarNever)
	t.SetSelectable(true, false)
	t.SetSortClicked(false)
	t.Select(0, 0)
	t.SetFixed(1, 1)

	trips := &trips{Table: t}
	trips.SetBorder(true)
	trips.setEntries(g)
	trips.setKeybinding(g)
	return trips
}

func (t *trips) name() string {
	return `trips`
}

func (t *trips) setKeybinding(g *Gui) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		switch event.Key() {
		case tcell.KeyEnter:
			g.state.navigate.update("detail")
			g.inspectTrip()
		case tcell.KeyCtrlR:
			t.setEntries(g)
		}

		switch event.Rune() {
		case 'n':
			g.createTripForm()
		case 'm':
			g.modifyTripForm()
		case 'd':
			g.deleteTrip()
		}

		return event
	})
}

func (t *trips) entries(g *Gui) {
	trips, err := g.reg.GetAllTrips()
	if err != nil {
		return
	}

	var filteredTrips []*model.Log
	for _, trip := range trips {
		if t.search(trip) {
			continue
		}
		filteredTrips = append(filteredTrips, trip)
	}
	g.state.resources.trips = filteredTrips
}

func (t *trips) setEntries(g *Gui) {
	t.entries(g)
	t.Clear()

	headers := [][]byte{
		[]byte("Date"),
		[]byte("Cave"),
		[]byte("Names"),
	}

	for i, header := range headers {
		t.SetCell(0, i, &cview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           cview.AlignLeft,
			Color:           cview.Styles.PrimaryTextColor,
			BackgroundColor: cview.Styles.PrimitiveBackgroundColor,
			Attributes:      tcell.AttrBold,
		})
	}

	var cell *cview.TableCell
	for i, trip := range g.state.resources.trips {
		cell = cview.NewTableCell(trip.Date)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		t.SetCell(i+1, 0, cell)

		cell = cview.NewTableCell(trip.Cave)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		t.SetCell(i+1, 1, cell)

		cell = cview.NewTableCell(trip.Names)
		cell.SetTextColor(cview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(0)
		cell.SetExpansion(2)
		t.SetCell(i+1, 2, cell)
	}

	t.Select(1, 1)
}

func (t *trips) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		t.setEntries(g)
	})
}

func (t *trips) focus(g *Gui) {
	t.SetSelectable(true, false)
	g.app.SetFocus(t)
}

func (t *trips) unfocus() {
	t.SetSelectable(false, false)
}

func (t *trips) setFilter(col, term, action string) {
	t.filterCol = col
	t.filterTerm = term
	t.filterAction = action // Unused in this window.
}

func (t *trips) monitoringTrips(g *Gui) {
	ticker := time.NewTicker(5 * time.Minute)

LOOP:
	for {
		select {
		case <-ticker.C:
			t.updateEntries(g)
		case <-g.state.stopChans["trips"]:
			ticker.Stop()
			break LOOP
		}
	}
}

func (t *trips) search(trip *model.Log) bool {
	// Below *looks* goofy but it all makes sense considering this funciton
	// needs to return false normally!!
	switch t.filterCol {
	case "cave", "":
		if strings.Index(strings.ToLower(trip.Cave), t.filterTerm) == -1 {
			return true
		}
		return false
	default:
		return false
	}
}
