package tui

import (
	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/pkg/model"
	"github.com/idlephysicist/cave-logger/internal/pkg/keeper"
)

type panels struct {
	currentPanel int
	panel        []panel
}

type resources struct {
	entries []*model.Entry
	caves   []*model.Cave
	cavers  []*model.Caver
}

type state struct {
	panels    panels
	navigate  *navigate
	resources resources
	stopChans map[string]chan int
}

type Tui struct {
	app   *tview.Application
	pages *tview.Pages
	state *state
	db    *keeper.Keeper
}



func New() *Tui {
	return &Tui{
		app: tview.NewApplication(),
		state: newState(),
	}
}

// Start start application
func (t *Tui) Start() error {
	t.initPanels()
	t.startMonitoring()
	if err := t.app.Run(); err != nil {
		t.app.Stop()
		return err
	}

	return nil
}

func (t *Tui) Stop() error {
	t.stopMonitoring()
	t.app.Stop()
	return nil
}

//
// -- Internal Functions -------------------------------------------------------
//

func (t *Tui) selectedEntry() *model.Entry {
	row, _ := t.entriesPanel().GetSelection()
	if len(t.state.resources.entries) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return t.state.resources.entries[row-1]
}


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
	return nil
}

func (t *Tui) initPanels() {
	entries := newEntries(t)

	t.state.panels.panel = append(t.state.panels.panel, entries)
	t.state.navigate = navi

	grid := tview.NewGrid().SetRows(2, 0, 0, 0, 0, 0, 2).
		AddItem(info, 0, 0, 1, 1, 0, 0, true).
		AddItem(entries, 1, 0, 1, 1, 0, 0, true).
		AddItem(navi, 6, 0, 1, 1, 0, 0, true)

	t.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	t.app.SetRoot(t.pages, true)
	t.switchPanel("entries")
}

func (t *Tui) switchPanel(panelName string) {
	for i, panel := range t.state.panels.panel {
		if panel.name() == panelName {
			t.state.navigate.update(panelName)
			panel.focus(g)
			t.state.panels.currentPanel = i
		} else {
			panel.unfocus()
		}
	}
}

func (t *Tui) closeAndSwitchPanel(removePanel, switchPanel string) {
	t.pages.RemovePage(removePanel).ShowPage("main")
	t.switchPanel(switchPanel)
}

func (t *Tui) modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}