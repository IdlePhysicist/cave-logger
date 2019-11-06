package gui

import (
	"fmt"

	"github.com/gdamore/tcell"
	//"github.com/rivo/tview"
	"github.com/idlephysicist/cave-logger/internal/model"
)

var inspectorFormat = map[string]string{
	`trips`   : "Date: %s\nCave: %s\nCavers: %s\nNotes: %s",
	`people`  : "First: %s\nLast: %s\nClub: %s\nCount: %d",
	`locations`: "Name: %s\nRegion: %s\nCountry: %s\nSRT: %v\nVisits: %d",
}

var inputWidth = 70

func (g *Gui) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Rune() {
	case 'l':
		g.goTo(`trips`)
	case 'o':
		g.goTo(`caves`)
	case 'p':
		g.goTo(`cavers`)
	case 'q':
		g.Stop()
	//case '/':
	//	g.filter()
	}
	
	/*switch event.Key() {
	case tcell.KeyTab:
		g.nextPage()
	case tcell.KeyBacktab:
		g.prevPage()
	case tcell.KeyRight:
		g.nextPage()
	case tcell.KeyLeft:
		g.prevPage()
	}*/
}

func (g *Gui) inspectTrip() {
	selected := g.selectedTrip()

	trip, err := g.db.GetLog(selected.ID)
	if err != nil {
		return
	}

	g.inspectorPanel().setEntry(g.formatTrip(trip))
}

func (g *Gui) inspectCave() {
	selected := g.selectedCave()

	cave, err := g.db.GetCave(selected.ID)
	if err != nil {
		return
	}

	g.inspectorPanel().setEntry(g.formatCave(cave))
}

func (g *Gui) inspectPerson() {
	selected := g.selectedPerson()

	caver, err := g.db.GetCaver(selected.ID)
	if err != nil {
		return
	}

	g.inspectorPanel().setEntry(g.formatPerson(caver))
}

//
// Formatting Functions
//
func (g *Gui) formatTrip(trip *model.Log) string {
	return fmt.Sprintf(inspectorFormat[`trips`], trip.Date, trip.Cave, trip.Names, trip.Notes)
}

func (g *Gui) formatCave(l *model.Cave) string {
	return fmt.Sprintf(inspectorFormat[`locations`], l.Name, l.Region, l.Country, l.SRT, l.Visits)
}

func (g *Gui) formatPerson(p *model.Caver) string {
	return fmt.Sprintf(inspectorFormat[`people`], p.First, p.Last, p.Club, p.Count)
}