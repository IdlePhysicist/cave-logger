package gui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/model"
)

var inspectorFormat = map[string]string{
	`trips` : "Date: %s\nCave: %s\nCavers: %s\nNotes: %s",
	`cavers`: "Name: %s\nClub: %s\nCount: %d",
	`caves` : "Name: %s\nRegion: %s\nCountry: %s\nSRT: %v\nVisits: %d",
}

func (g *Gui) setGlobalKeybinding(event *tcell.EventKey) {
	/*switch event.Key() {
	case tcell.KeyTAB:
		g.nextPage()
	}*/

	switch event.Rune() {
	case 'q':
		g.Stop()
	case '1':
		g.goTo("trips")
	case '2':
		g.goTo("cavers")
	case '3':
		g.goTo("caves")
	case '/':
		g.filter()
	}
}

func (g *Gui) filter() {
	currentPanel := g.state.panels.panel[g.state.panels.currentPanel]
	currentPanel.setFilterWord("")
	currentPanel.updateEntries(g)

	viewName := "filter"
	searchInput := tview.NewInputField().SetLabel("Parameter")
	searchInput.SetLabelWidth(10)
	searchInput.SetTitle(" Filter ")
	searchInput.SetTitleAlign(tview.AlignLeft)
	searchInput.SetBorder(true)

	closeSearchInput := func() {
		g.closeAndSwitchPanel(viewName, g.state.panels.panel[g.state.panels.currentPanel].name())
	}

	searchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			closeSearchInput()
		}
	})

	searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			closeSearchInput()
		}
		return event
	})

	searchInput.SetChangedFunc(func(text string) {
		currentPanel.setFilterWord(text)
		currentPanel.updateEntries(g)
	})

	g.pages.AddAndSwitchToPage(viewName, g.modal(searchInput, 80, 3), true).ShowPage("main")
}

//
// INSPECTION FUNCS
//

func (g *Gui) inspectTrip() {
	selected := g.selectedTrip()

	trip, err := g.db.GetTrip(selected.ID)
	if err != nil {
		return
	}

	g.inspectorPanel().setEntry(g.formatTrip(trip))
}

func (g *Gui) inspectCave() {
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
}

//
// Formatting Functions
//
func (g *Gui) formatTrip(trip *model.Log) string {
	return fmt.Sprintf(inspectorFormat[`trips`], trip.Date, trip.Cave, trip.Names, trip.Notes)
}

func (g *Gui) formatCave(l *model.Cave) string {
	return fmt.Sprintf(inspectorFormat[`caves`], l.Name, l.Region, l.Country, l.SRT, l.Visits)
}

func (g *Gui) formatPerson(p *model.Caver) string {
	return fmt.Sprintf(inspectorFormat[`cavers`], p.Name, p.Club, p.Count)
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

/*
func (g *Gui) nextPage() {
	slide, _ := strconv.Atoi(g.state.tabBar.GetHighlights()[0])
	slide = (slide + 1) % g.pages.GetPageCount()
	//g.state.tabBar.Highlight(strconv.Itoa(slide)).ScrollToHighlight()
	g.goTo(g.selectPage(slide - 1, 0)) // NOTE: If the Highlight func is fixed for the tab bar then this line will not be required
}
*/
