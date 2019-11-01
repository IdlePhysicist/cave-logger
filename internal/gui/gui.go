package gui

import (
	//"context"

	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/db"
	//"github.com/idlephysicist/cave-logger/internal/model"
)

type pages struct {
	currentPage int
	page 				[]panel
}

type resources struct {
	trips  []*trips
	cavers []*cavers
	//caves  []*caves
}

type state struct {
	pages 	 	pages
	//navigate 	*navigate
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
	g.initPages()
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
	for _, page := range g.state.pages.page {
		if page.name() == `trips` {
			return page.(*trips)
		}
	}
	return nil
}

func (g *Gui) cavesPanel() *caves {
	for _, page := range g.state.pages.page {
		if page.name() == `caves` {
			return page.(*caves)
		}
	}
	return nil
}

func (g *Gui) caversPanel() *cavers {
	for _, page := range g.state.pages.page {
		if page.name() == `cavers` {
			return page.(*cavers)
		}
	}
	return nil
}

func (g *Gui) initPages() {
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	trips := newTrips(g)
	cavers := newCavers(g)
	caves  := newCaves(g)
	//navi  := newNavigate(g)
	//help  := newHelp(g)

	g.pages.AddPage(`trips`, trips, true, true)
	g.pages.AddPage(`cavers`, cavers, true, true)
	g.pages.AddPage(`caves`, caves, true, true)
	

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(g.pages, 0, 1, true).
		AddItem(info, 1, 1, false)


	g.app.SetRoot(layout, true)
	g.goTo(`trips`)
}

func (g *Gui) goTo(page string) {
	g.pages.SwitchToPage(page)
}
