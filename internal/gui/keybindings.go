package gui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/idlephysicist/cave-logger/internal/model"
)

var inspectorFormat = map[string]string{
	`trips`    : "Date: %s\nCave: %s\nCavers: %s\nNotes: %s",
	`people`   : "Name: %s\nClub: %s\nCount: %d",
	`locations`: "Name: %s\nRegion: %s\nCountry: %s\nSRT: %v\nVisits: %d",
}

var inputWidth = 70

func (g *Gui) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Rune() {
	case 'q':
		g.Stop()
	}
}

//
// INSPECTION FUNCS
//

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
	return fmt.Sprintf(inspectorFormat[`people`], p.Name, p.Club, p.Count)
}

//
// MISC
//

func (g *Gui) selectPage(row, col int) string {
	var p string
	switch row {
	case 0:
		p = `trips`
	case 1:
		p = `cavers`
	case 2:
		p = `caves`
	}
	return p
}