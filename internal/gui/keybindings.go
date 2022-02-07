package gui

import (
	"strings"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

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
	currentPanel.setFilter("", "")
	currentPanel.updateEntries(g)

	viewName := "filter"
	searchInput := cview.NewInputField()
	searchInput.SetLabel("Column/Parameter")
	searchInput.SetLabelWidth(17)
	searchInput.SetTitle(" Filter ")
	searchInput.SetTitleAlign(cview.AlignLeft)
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

	g.pages.AddAndSwitchToPage(viewName, g.modal(searchInput, 80, 3), true)
	g.pages.ShowPage("main")
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
