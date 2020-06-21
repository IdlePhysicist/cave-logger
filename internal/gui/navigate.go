package gui

import "github.com/rivo/tview"

type navigate struct {
	*tview.TextView
	keybindings map[string]string
}

func newNavigate() *navigate {
	navi := &navigate{
		TextView: tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor),
		keybindings: map[string]string{
			"trips" : " n: New Log Entry, m: Modify Log,  d: Remove Log, /: Filter, Enter: Inspect Detail ",
			"caves" : " n: New Cave, m: Modify Cave, d: Remove Cave, /: Filter ",
			"cavers": " n: New Caver, m: Modify Caver, d: Remove Caver, /: Filter ",
			"detail": " q | ESC: Exit Detail ",
		},
	}

	navi.SetBorder(false)
	return navi
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
