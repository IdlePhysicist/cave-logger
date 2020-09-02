package gui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	tview "gitlab.com/tslocum/cview"
)

func (g *Gui) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyRight:
		g.nextCol()
	case tcell.KeyLeft:
		g.prevCol()
	}

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
	case 'l':
		g.nextCol()
	case 'h':
		g.prevCol()
	}
}

func (g *Gui) filter() {
	currentPanel := g.state.panels.panel[g.state.panels.currentPanel]
	currentPanel.setFilter("", "")
	currentPanel.updateEntries(g)

	viewName := "filter"
	searchInput := tview.NewInputField().SetLabel("Column/Parameter")
	searchInput.SetLabelWidth(17)
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
		if strings.Contains(text, "/") {
			textSl := strings.Split(strings.ToLower(text), "/")

			if len(textSl) == 2 {
				currentPanel.setFilter(textSl[0], textSl[1])
				currentPanel.updateEntries(g)
			}
		}
	})

	g.pages.AddAndSwitchToPage(viewName, g.modal(searchInput, 80, 3), true).ShowPage("main")
}

//
// MISC
//

func (g *Gui) selectPage(row, col int) string {
	var p string // REVIEW: can this be improved?
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


func (g *Gui) nextCol() {
	currentPanel := g.state.panels.panel[g.state.panels.currentPanel]
	newSortedCol := (currentPanel.getSortedCol() + 1) % currentPanel.getColumnCount()

	currentPanel.setSortedCol(newSortedCol)
}


func (g *Gui) prevCol() {
	currentPanel := g.state.panels.panel[g.state.panels.currentPanel]
	newSortedCol := (currentPanel.getSortedCol() - 1) % currentPanel.getColumnCount()

	currentPanel.setSortedCol(newSortedCol)
}
