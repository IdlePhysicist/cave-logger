package tui

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
			"entries":     " p: pull image, i: import image, s: save image, Ctrl+l: load image, f: search image, /: filter d: remove image,\n c: create container, Enter: inspect image, Ctrl+r: refresh images list",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}