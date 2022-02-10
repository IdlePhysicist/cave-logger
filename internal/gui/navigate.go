package gui

import "code.rocketnine.space/tslocum/cview"

type navigate struct {
	*cview.TextView
	keybindings map[string]string
}

func newNavigate() *navigate {
	tv := cview.NewTextView()
	tv.SetTextColor(cview.Styles.PrimaryTextColor)

	navi := &navigate{
		TextView: tv,
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
