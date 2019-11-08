package gui

import (
	"github.com/rivo/tview"
)

var inputWidth = 70

//
// CREATE FUNCS
func (g *Gui) createTripForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Add Log Entry")
	form.SetTitleAlign(tview.AlignLeft)

	form.
		AddInputField("Date", "", inputWidth, nil, nil).
		AddInputField("Cave", "", inputWidth, nil, nil).
		AddInputField("Names", "", inputWidth, nil, nil).
		AddInputField("Notes", "", inputWidth, nil, nil).
		AddButton("Add", func() {
			g.createTrip(form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "trips")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)//.ShowPage("main")
	//REVIEW: main or trips ? ^^
}

func (g *Gui) createTrip(form *tview.Form) {
	err := g.db.AddLog(
		form.GetFormItemByLabel("Date").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Cave").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Names").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Notes").(*tview.InputField).GetText(),
	)
	if err != nil { // NOTE: Needs fixing
		g.warning(err.Error(), `OK`, `form`, func() {return})
		return
	}

	g.closeAndSwitchPanel(`form`, `trips`)
	g.app.QueueUpdateDraw(func() {
		g.tripsPanel().setEntries(g)
	})
}

//
// UPDATE FUNCS

//
// DELETE FUNCS
func (g *Gui) deleteTrip() {}
