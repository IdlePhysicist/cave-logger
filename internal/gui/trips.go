package gui

import (
	"strings"
	"time"

	"github.com/gdamore/tcell"
	tview "gitlab.com/tslocum/cview"

	"github.com/IdlePhysicist/cave-logger/internal/model"
)

type trips struct {
	*tview.Table
	trips			 chan *model.Log
	filterWord string
}

func newTrips(g *Gui) *trips {
	trips := &trips{
		Table: tview.NewTable().SetSelectable(true, false).Select(0,0).SetFixed(1,1),
		trips: make(chan *model.Log),
	}

	trips.SetTitle(``).SetTitleAlign(tview.AlignLeft)
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
	trips, err := g.db.GetAllTrips()
	if err != nil {
		return
	}

	var filteredTrips []*model.Log
	for _, trip := range trips {
		if strings.Index(trip.Cave, t.filterWord) == -1 {
			continue
		}
		filteredTrips = append(filteredTrips, trip)
	}
	g.state.resources.trips = filteredTrips
}

func (t *trips) setEntries(g *Gui) {
	t.entries(g)
	table := t.Clear()

	headers := []string{
		"Date",
		"Cave",
		"Names",
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

	for i, trip := range g.state.resources.trips {
		table.SetCell(i+1, 0, tview.NewTableCell(trip.Date).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetMaxWidth(30).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(trip.Cave).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetMaxWidth(30).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(trip.Names).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetMaxWidth(0).
			SetExpansion(2))
	}
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

func (t *trips) setFilterWord(word string) {
	t.filterWord = word
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
