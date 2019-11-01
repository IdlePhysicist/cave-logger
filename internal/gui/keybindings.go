package gui

import (
	//"context"
//	"fmt"


	"github.com/gdamore/tcell"
	//"github.com/rivo/tview"
//	"github.com/skanehira/docui/common"
)

var inputWidth = 70

func (g *Gui) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Rune() {
	case 'l':
		g.goTo(`trips`)
	case 'k':
		g.goTo(`caves`)
	case 'j':
		g.goTo(`cavers`)
	case 'q':
		g.Stop()
	//case '/':
	//	g.filter()
	}
	
	/*
	switch event.Key() {
	case tcell.KeyTab:
		g.nextPanel()
	case tcell.KeyBacktab:
		g.prevPanel()
	case tcell.KeyRight:
		g.nextPanel()
	case tcell.KeyLeft:
		g.prevPanel()
	}*/
}