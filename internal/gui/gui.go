package gui

import (
	"fmt"
	"strconv"
	"strings"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"

	"github.com/idlephysicist/cave-logger/internal/model"
	"github.com/idlephysicist/cave-logger/internal/register"
)

type panels struct {
	currentPanel int
	panel        []panel
}

type resources struct {
	trips     []*model.Log
	people    []*model.Caver
	locations []*model.Cave
	//statsLocations []*model.Statistic
}

type state struct {
	panels    panels
	navigate  *navigate
	tabBar    *cview.TextView
	resources resources
	stopChans map[string]chan int
}

func newState() *state {
	return &state{
		stopChans: make(map[string]chan int),
	}
}

type Gui struct {
	app   *cview.Application
	pages *cview.Pages
	state *state
	reg   *register.Register
	//statsLocations *statsLocations
}

func New(reg *register.Register) *Gui {
	return &Gui{
		app:   cview.NewApplication(),
		pages: cview.NewPages(),
		state: newState(),
		reg:   reg,
	}
}

func (g *Gui) ProcessColors(colors map[string]string) {
	for color, hex := range colors {
		if hex == "" {
			continue
		}
		switch color {
		case "primitiveBackground":
			cview.Styles.PrimitiveBackgroundColor = tcell.GetColor(hex)
		case "contrastBackground":
			cview.Styles.ContrastBackgroundColor = tcell.GetColor(hex)
		case "moreContrastBackground":
			cview.Styles.MoreContrastBackgroundColor = tcell.GetColor(hex)
		case "border":
			cview.Styles.BorderColor = tcell.GetColor(hex)
		case "title":
			cview.Styles.TitleColor = tcell.GetColor(hex)
		case "graphics":
			cview.Styles.GraphicsColor = tcell.GetColor(hex)
		case "primaryText":
			cview.Styles.PrimaryTextColor = tcell.GetColor(hex)
		case "secondaryText":
			cview.Styles.SecondaryTextColor = tcell.GetColor(hex)
		case "tertiaryText":
			cview.Styles.TertiaryTextColor = tcell.GetColor(hex)
		case "inverseText":
			cview.Styles.InverseTextColor = tcell.GetColor(hex)
		case "contrastSecondaryText":
			cview.Styles.ContrastSecondaryTextColor = tcell.GetColor(hex)
		}
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

/*func (g *Gui) statsLocationsPanel() *statsLocations {
	for _, panel := range g.state.panels.panel {
		if panel.name() == `statsLocations` {
			return panel.(*statsLocations)
		}
	}
	return nil
}
*/

func (g *Gui) initPanels() {

	g.state.tabBar = newTabBar(g)

	// Page definitions
	trips := newTrips(g)
	cavers := newCavers(g)
	caves := newCaves(g)

	/*
		// NOTE: I would really like to get this working as it would be far neater.
		// The issue is with the three pages being of different types.
		// cannot use pg (type panel) as type cview.Primitive in argument to g.pages.AddPage:
		// panel does not implement cview.Primitive (missing Blur method)
		for idx, pg := range []panel{trips, cavers, caves} {
			name := pg.name()
			g.pages.AddPage(name, pg, true, idx == 0)
			fmt.Fprintf(g.state.tabBar, ` %d ["%d"][darkcyan]%s[white][""]  `, idx+1, idx, strings.Title(name))
		}
		g.state.tabBar.Highlight("0")
	*/

	// Add pages to the "book"
	g.pages.AddPage(`trips`, trips, true, true)
	fmt.Fprintf(g.state.tabBar, `  ["%d"]%d %s[""] `, 0, 1, strings.Title(trips.name()))
	g.pages.AddPage(`cavers`, cavers, true, true)
	fmt.Fprintf(g.state.tabBar, `  ["%d"]%d %s[""] `, 1, 2, strings.Title(cavers.name()))
	g.pages.AddPage(`caves`, caves, true, true)
	fmt.Fprintf(g.state.tabBar, `  ["%d"]%d %s[""] `, 2, 3, strings.Title(caves.name()))

	g.state.tabBar.Highlight("0")

	// Panels
	statusBar := newNavigate()

	g.state.panels.panel = append(g.state.panels.panel, trips)
	g.state.panels.panel = append(g.state.panels.panel, cavers)
	g.state.panels.panel = append(g.state.panels.panel, caves)

	g.state.navigate = statusBar

	// Arange the windows / tiles
	layout := cview.NewFlex()
	layout.SetDirection(cview.FlexRow)
	layout.AddItem(g.state.tabBar, 1, 1, false)
	layout.AddItem(g.pages, 0, 16, true)
	layout.AddItem(statusBar, 1, 1, false)

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
			g.state.tabBar.Highlight(strconv.Itoa(i))
			g.state.tabBar.ScrollToHighlight()
		} else {
			panel.unfocus()
		}
	}
}

func (g *Gui) closeAndSwitchPanel(removePanel, switchTo string) {
	g.pages.RemovePage(removePanel)
	g.pages.ShowPage("main")
	num := 0
	switch switchTo {
	case `cavers`:
		num = 1
	case `caves`:
		num = 2
	default:
		num = 0
	}
	g.goTo(g.selectPage(num, 0))
}

func (g *Gui) currentPage() int {
	return g.state.panels.currentPanel
}

func (g *Gui) modal(p cview.Primitive, width, height int) cview.Primitive {
	grid := cview.NewGrid()
	grid.SetColumns(0, width, 0)
	grid.SetRows(0, height, 0)
	grid.AddItem(p, 1, 1, 1, 1, 0, 0, true)
	return grid
}

func (g *Gui) warning(message, page string, labels []string, doneFunc func()) {
	modal := cview.NewModal()
	modal.SetText(message)
	modal.AddButtons(labels)
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		g.closeAndSwitchPanel("modal", page)
		if buttonLabel == labels[0] {
			doneFunc()
		}
	})

	g.pages.AddAndSwitchToPage("modal", g.modal(modal, 80, 29), true)
}

//
// Functions for returning the selected item in the table
// REVIEW: There might be better ways of doing this.
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

func (g *Gui) selectedLocation() *model.Cave {
	row, _ := g.cavesPanel().GetSelection()
	if len(g.state.resources.locations) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.state.resources.locations[row-1]
}

func (g *Gui) selectedPerson() *model.Caver {
	row, _ := g.caversPanel().GetSelection()
	if len(g.state.resources.people) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.state.resources.people[row-1]
}
