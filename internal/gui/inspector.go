package gui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	tview "gitlab.com/tslocum/cview"

	"github.com/idlephysicist/cave-logger/internal/model"
)

var inspectorFormat = map[string]string{
	`trips`:  "\tDate: %s\n\tCave: %s\n\tCavers: %s\n\tNotes: %s",
	`cavers`: "\tName: %s\n\tClub: %s\n\tCount: %d\n\tLast Trip: %s\n\tNotes: %s",
	`caves`:  "\tName: %s\n\tRegion: %s\n\tCountry: %s\n\tSRT: %s\n\tVisits: %d\n\tLast Visit: %s\n\tNotes: %s",
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
		g.warning("No trips in table", `trips`, []string{`OK`}, func() { return })
		return
	}

	trip, err := g.reg.GetTrip(selected.ID)
	if err != nil {
		return
	}

	g.displayInspect(g.formatTrip(trip), "trips")
}

func (g *Gui) inspectCave() {
	selected := g.selectedLocation()

	cave, err := g.reg.GetCave(selected.ID)
	if err != nil {
		return
	}

	g.displayInspect(g.formatCave(cave), "caves")
}

func (g *Gui) inspectCaver() {
	selected := g.selectedPerson()

	caver, err := g.reg.GetCaver(selected.ID)
	if err != nil {
		return
	}

	g.displayInspect(g.formatPerson(caver), "cavers")
}

//
// Formatting Functions
//
func (g *Gui) formatTrip(t *model.Log) string {
	return fmt.Sprintf(inspectorFormat[`trips`], t.Date, t.Cave, t.Names, t.Notes)
}

func (g *Gui) formatCave(c *model.Cave) string {
	return fmt.Sprintf(inspectorFormat[`caves`],
		c.Name, c.Region, c.Country, yesOrNo(c.SRT), c.Visits, c.LastVisit, c.Notes,
	)
}

func (g *Gui) formatPerson(c *model.Caver) string {
	return fmt.Sprintf(inspectorFormat[`cavers`], c.Name, c.Club, c.Count, c.LastTrip, c.Notes)
}
