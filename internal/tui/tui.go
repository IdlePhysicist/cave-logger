package tui

import (
	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/pkg/model"
)

type panels struct {
	currentPanel int
	panel        []panel
}

type resources struct {
	logs   []*model.Entry
	caves  []*model.Caver
	cavers []*model.Cave
}

type state struct {
	panels    panels
	navigate  *navigate
	resources resources
	stopChans map[string]chan int
}

type Tui struct {
	app *tview.Application
	pages *tview.Pages
	state *state
}



func New() *Tui {
	return &Tui{
		app: tview.NewApplication(),
		state: newState(),
	}
}

func (t *Tui) Start() error {



	
	return nil
}


//
// -- Internal Functions -------------------------------------------------------
//


func newState() *state {
	return &state{
		stopChans: make(map[string]chan int),
	}
}

func (t *Tui) entriesPanel() *entries {
	for _, panel := range t.state.panels.panel {
		if panel.name() == `entries` {
			return panel.(*entries)
		}
	}
}