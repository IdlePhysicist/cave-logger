package tui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	//"github.com/idlephysicist/cave-logger/internal/pkg/model"
	"github.com/idlephysicist/cave-logger/internal/pkg/keeper"
)

type entries struct {
	*tview.Table
	filterWord string
	db *keeper.Keeper
}

func newEntries(k *keeper.Keeper, t *Tui) *entries {
	entries := &entries{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		db: k,
	}

	entries.SetTitle(` Logs `).SetTitleAlign(tview.AlignLeft)
	entries.SetBorder(true)
	entries.setContents(t)
	entries.setKeybinding(t)
	return entries
}

func (e *entries) name() string {
	return `entries`
}

func (e *entries) setContents(t *Tui) {
	e.contents(t)
	table := e.Clear()

	headers := []string{
		`Date`,
		`Cave`,
		`Names`,
		`Notes`,
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text: header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold,
		})
	}

	for i, entry := range t.state.resources.entries {
		table.SetCell(i+1, 0, tview.NewTableCell(entry.Date).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(entry.Cave).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(entry.Names).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 3, tview.NewTableCell(entry.Notes).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))
		// rm the ntoes one?
	}

}

func (e *entries) setKeybinding(t *Tui) {
	e.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			t.inspectEntry()
		}

		switch event.Rune() {
		case 'n':
			t.addEntryForm()
		case 'r':
			t.removeEntry()
		//case 'f':
		//	t.searchInputField()
		}

		return event
	})
}

func (e *entries) contents(t *Tui) {
	rows, err := e.db.QueryLogs(`-1`)
	if err != nil {
		return
	}

	t.state.resources.entries = rows //make([]*model.Entry, 0)

	/*for _, row := range rows {
		t.state.resources.entries = append(t.state.resources.entries, &model.Entry)
	}*/
}

// --------------------

func (e *entries) updateContents(t *Tui) {
	t.app.QueueUpdateDraw(func() {
		e.setContents(t)
	})
}

func (e *entries) focus(t *Tui) {
	e.SetSelectable(true, false)
	t.app.SetFocus(e)
}

func (e *entries) unfocus() {
	e.SetSelectable(false, false)
}

func (e *entries) setFilterWord(word string) {
	e.filterWord = word
}

func (e *entries) monitoringImages(t *Tui) {
	//common.Logger.Info("start monitoring images")
	ticker := time.NewTicker(5 * time.Second)

	LOOP:
		for {
			select {
			case <-ticker.C:
				e.updateContents(t)
			case <-t.state.stopChans["image"]:
				ticker.Stop()
				break LOOP
			}
		}
		//common.Logger.Info("stop monitoring images")
}