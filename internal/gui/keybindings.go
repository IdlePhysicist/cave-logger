package gui

import (
	//"context"
//	"fmt"
	//"strconv"


	"github.com/gdamore/tcell"
	//"github.com/rivo/tview"
//	"github.com/skanehira/docui/common"
)

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

func (g *Gui) nextPage() {
	//idx := (g.state.pages.currentPage + 1) % len(g.state.pages.page)
	//g.switchPage(g.state.pages.page[idx].name())
	currentSlide := (g.state.pages.currentPage + 1) % 3//len(slides)
	//info.Highlight(strconv.Itoa(currentSlide)).ScrollToHighlight()
	g.pages.SwitchToPage(g.state.pages.page[currentSlide].name())
}

func (g *Gui) prevPage() {
	g.state.pages.currentPage--

	if g.state.pages.currentPage < 0 {
		g.state.pages.currentPage = len(g.state.pages.page) - 1
	}

	idx := (g.state.pages.currentPage) % len(g.state.pages.page)
	g.switchPage(g.state.pages.page[idx].name())
}

func (g *Gui) switchPage(pageName string) {
	for i, page := range g.state.pages.page {
		if page.name() == pageName {
			//g.state.navigate.update(pageName)
			page.focus(g)
			g.state.pages.currentPage = i
		} else {
			page.unfocus()
		}
	}
}
