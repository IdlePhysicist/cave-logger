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
  navi := &navigate{
    TextView: tview.NewTextView().SetTextColor(tcell.ColorYellow),
    keybindings: map[string]string{
      "trips": " n: New Log Entry, m: Modify Log,  d: Remove Log, /: filter, Enter: inspect ",
      "locations": " n: New Cave, m: Modify Cave, d: Remove Cave, /: filter, Enter: Inspect ",
      "people": " n: New Caver, m: Modify Caver, d: Remove Caver, /: filter, Enter: Inspect ",
    },
  }

  navi.SetBorder(true)
  return navi
}

func (n *navigate) update(panel string) {
  n.SetText(n.keybindings[panel])
}
