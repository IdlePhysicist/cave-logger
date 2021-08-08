package gui

import tview "gitlab.com/tslocum/cview"

type navigate struct {
	*tview.TextView
	keybindings map[string]string
}

func newNavigate() *navigate {
	navi := &navigate{
		TextView: tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor),
		keybindings: map[string]string{
			"trips":  " n: New Trip, m: Modify Trip,  d: Remove Trip, /: Filter, Enter: Inspect Detail ",
			"caves":  " n: New Cave, m: Modify Cave, d: Remove Cave, /: Filter, Enter: Inspect Detail ",
			"cavers": " n: New Caver, m: Modify Caver, d: Remove Caver, /: Filter, Enter: Inspect Detail ",
			"detail": " q | ESC: Exit Detail ",
		},
	}

	navi.SetBorder(false)
	return navi
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
