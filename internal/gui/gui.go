package gui

import (
	//"context"
	//"fmt"

	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/db"
	"github.com/idlephysicist/cave-logger/internal/model"
)

type panels struct {
	currentPanel int
	panel 			 []panel
}

type resources struct {
	trips  []*model.Log //trips
	cavers []*model.Caver
	caves  []*model.Cave
	menu 	 []string
	statsLocations  []*model.Statistic
	statsPeople  		[]*model.Statistic
	timeWindow      []*model.Statistic
}

type state struct {
	panels 	 	panels
	//insp      *inspector
	navigate *navigate
	resources resources
	stopChans map[string]chan int
}

func newState() *state {
	return &state{
		stopChans: make(map[string]chan int),
	}
}

type Gui struct {
	app 	*tview.Application
	pages *tview.Pages
	state *state
	db    *db.Database
	statsLocations *statsLocations
	statsPeople 	 *statsPeople
}

func New(db *db.Database) *Gui {
	return &Gui{
		app: tview.NewApplication(),
		pages: tview.NewPages(),
		state: newState(),
		db: db,
	}
}

// Start start application
func (g *Gui) Start() error {
	g.initPanels()
	g.startMonitoring()
	if err := g.app.Run(); err != nil {
		g.app.Stop()
		return err
	}

	return nil
}

func (g *Gui) Stop() {
	g.stopMonitoring()
	g.app.Stop()
}

// Page "definitions"

func (g *Gui) tripsPanel() *trips {
	for _, panel := range g.state.panels.panel {
		if panel.name() == `trips` {
			return panel.(*trips)
		}
	}
	return nil
}

func (g *Gui) cavesPanel() *caves {
	for _, panel := range g.state.panels.panel {
		if panel.name() == `caves` {
			return panel.(*caves)
		}
	}
	return nil
}

func (g *Gui) caversPanel() *cavers {
	for _, panel := range g.state.panels.panel {
		if panel.name() == `cavers` {
			return panel.(*cavers)
		}
	}
	return nil
}

func (g *Gui) inspectorPanel() *inspector {
	for _, panel := range g.state.panels.panel {
		if panel.name() == `inspector` {
			return panel.(*inspector)
		}
	}
	return nil
}

func (g *Gui) statsLocationsPanel() *statsLocations {
	for _, panel := range g.state.panels.panel {
		if panel.name() == `statsLocations` {
			return panel.(*statsLocations)
		}
	}
	return nil
}

func (g *Gui) statsPeoplePanel() *statsPeople {
	for _, panel := range g.state.panels.panel {
		if panel.name() == `statsPeople` {
			return panel.(*statsPeople)
		}
	}
	return nil
}


func (g *Gui) initPanels() {
	// Page definitions
	trips  := newTrips(g)
	cavers := newCavers(g)
	caves  := newCaves(g)

	// Add pages to the "book"
	g.pages.AddPage(`trips`, trips, true, true)
	g.pages.AddPage(`cavers`, cavers, true, true)
	g.pages.AddPage(`caves`, caves, true, true)
	
	// Panels
	menu := newMenu(g)
	statsPeople := newStatsPeople(g)
	statsLocations := newStatsLocations(g)
	timeWindow := newTimeWindow(g)
	inspector := newInspector(g)
	navigate := newNavigate()

	g.state.panels.panel = append(g.state.panels.panel, trips)
	g.state.panels.panel = append(g.state.panels.panel, cavers)
	g.state.panels.panel = append(g.state.panels.panel, caves)
	g.state.panels.panel = append(g.state.panels.panel, menu)
	g.state.panels.panel = append(g.state.panels.panel, statsPeople)
	g.state.panels.panel = append(g.state.panels.panel, statsLocations)
	g.state.panels.panel = append(g.state.panels.panel, timeWindow)
	g.state.panels.panel = append(g.state.panels.panel, inspector)

	g.state.navigate = navigate

	// Arange the windows / tiles
	layout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(menu, 0, 5, false).
			AddItem(statsPeople, 0, 20, false).
			AddItem(statsLocations, 0, 20, false).
			AddItem(timeWindow, 0, 2, false),
			0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(g.pages, 0, 16, true).
			AddItem(inspector, 0, 8, false).
			AddItem(navigate, 0, 1, false),
			0, 6, true)

	g.statsPeople = statsPeople
	g.statsLocations = statsLocations

	g.app.SetRoot(layout, true)
	g.goTo(`trips`)
}

func (g *Gui) goTo(page string) {
	g.pages.SwitchToPage(page)
	g.switchPanel(page)
}

func (g *Gui) switchPanel(panelName string) {
	for i, panel := range g.state.panels.panel {
		if panel.name() == panelName {
			g.state.navigate.update(panelName)
			panel.focus(g)
			g.state.panels.currentPanel = i
		} else {
			panel.unfocus()
		}
	}
}

func (g *Gui) closeAndSwitchPanel(removePanel, switchPanel string) {
	g.pages.RemovePage(removePanel).ShowPage("main")
	g.switchPanel(switchPanel)
}

func (g *Gui) currentPage() int {
	return g.state.panels.currentPanel
}

func (g *Gui) modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}


func (g *Gui) selectedTrip() *model.Log {
	row, _ := g.tripsPanel().GetSelection()
	if len(g.state.resources.trips) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.state.resources.trips[row-1]
}

func (g *Gui) selectedCave() *model.Cave {
	row, _ := g.cavesPanel().GetSelection()
	if len(g.state.resources.caves) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.state.resources.caves[row-1]
}

func (g *Gui) selectedPerson() *model.Caver {
	row, _ := g.caversPanel().GetSelection()
	if len(g.state.resources.cavers) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.state.resources.cavers[row-1]
}
