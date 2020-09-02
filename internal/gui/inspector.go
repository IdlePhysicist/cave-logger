package gui

import (
	"fmt"

	tview "gitlab.com/tslocum/cview"
	"github.com/gdamore/tcell/v2"

	"github.com/idlephysicist/cave-logger/internal/model"
)

var inspectorFormat = map[string]string{
	`trips` : "\tDate: %s\n\tCave: %s\n\tCavers: %s\n\tNotes: %s",
	`cavers`: "\tName: %s\n\tClub: %s\n\tCount: %d",
	`caves` : "\tName: %s\n\tRegion: %s\n\tCountry: %s\n\tSRT: %v\n\tVisits: %d",
}

func (g *Gui) displayInspect(data, page string) {
	text := tview.NewTextView()
	text.SetTitle(" Detail ").SetTitleAlign(tview.AlignLeft)
	text.SetBorder(true)
	text.SetText(data)

	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			g.closeAndSwitchPanel("detail", page)
		}
		return event
	})

	g.pages.AddAndSwitchToPage("detail", text, true)
}

//
// INSPECTION FUNCS
//

func (g *Gui) inspectTrip() {
	selected := g.selectedTrip()

	if selected == nil {
		g.warning("No trips in table", `trips`, []string{`OK`}, func() {return})
		return
	}

	trip, err := g.db.GetTrip(selected.ID)
	if err != nil {
		return
	}

	g.displayInspect(g.formatTrip(trip), "trips")
}

/*func (g *Gui) inspectCave() {
	selected := g.selectedLocation()

	cave, err := g.db.GetLocation(selected.ID)
	if err != nil {
		return
	}

	g.inspectorPanel().setEntry(g.formatCave(cave))
}

func (g *Gui) inspectPerson() {
	selected := g.selectedPerson()

	caver, err := g.db.GetPerson(selected.ID)
	if err != nil {
		return
	}

	g.inspectorPanel().setEntry(g.formatPerson(caver))
}*/

//
// Formatting Functions
//
func (g *Gui) formatTrip(trip *model.Log) string {
	return fmt.Sprintf(inspectorFormat[`trips`], trip.Date, trip.Cave, trip.Names, trip.Notes)
}

/*func (g *Gui) formatCave(l *model.Cave) string {
	return fmt.Sprintf(inspectorFormat[`caves`], l.Name, l.Region, l.Country, l.SRT, l.Visits)
}

func (g *Gui) formatPerson(p *model.Caver) string {
	return fmt.Sprintf(inspectorFormat[`cavers`], p.Name, p.Club, p.Count)
}*/

