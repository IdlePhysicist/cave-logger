package gui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type navigate struct {
	*tview.TextView
	keybindings map[string]string
}

func newNavigate() *navigate {
	return &navigate{
		TextView: tview.NewTextView().SetTextColor(tcell.ColorYellow),
		keybindings: map[string]string{
			"trips": " n: New Log Entry , m: Modify Log,  d: Remove Log\n /: filter, Enter: inspect ",
			"caves": " n: New Cave , m: Modify Cave, d: Remove Cave\n /: filter, Enter: Inspect ",
			"cavers": " n: New Caver , m: Modify Caver, d: Remove Caver\n /: filter, Enter: Inspect ",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
