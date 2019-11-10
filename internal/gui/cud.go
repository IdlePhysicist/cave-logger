package gui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

var inputWidth = 70

//
// CREATE FUNCS
func (g *Gui) createTripForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Add Trip ")
	form.SetTitleAlign(tview.AlignLeft)

	form.
		AddInputField("Date", time.Now().Format(`2006-01-02`), inputWidth, nil, nil).
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
	err := g.db.AddTrip(
		form.GetFormItemByLabel("Date").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Cave").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Names").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Notes").(*tview.InputField).GetText(),
	)
	if err != nil { // NOTE: Needs fixing
		g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
		return
	}

	g.closeAndSwitchPanel(`form`, `trips`)
	g.app.QueueUpdateDraw(func() {
		g.tripsPanel().setEntries(g)
	})
}

func (g *Gui) createLocationForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Add Location ")
	form.SetTitleAlign(tview.AlignLeft)

	form.
		AddInputField("Name", "", inputWidth, nil, nil).
		AddInputField("Region", "", inputWidth, nil, nil).
		AddInputField("Country", "", inputWidth, nil, nil).
		AddCheckbox("SRT", false, nil).
		AddButton("Add", func() {
			g.createLocation(form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "locations")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)//.ShowPage("main")
	//REVIEW: main or trips ? ^^
}

func (g *Gui) createLocation(form *tview.Form) {
	err := g.db.AddLocation(
		form.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Region").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Country").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("SRT").(*tview.Checkbox).IsChecked(),
	)
	if err != nil { // NOTE: Needs fixing
		g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
		return
	}

	g.closeAndSwitchPanel(`form`, `locations`)
	g.app.QueueUpdateDraw(func() {
		g.locationsPanel().setEntries(g)
	})
}

func (g *Gui) createPersonForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Add Person ")
	form.SetTitleAlign(tview.AlignLeft)

	form.
		AddInputField("Name", "", inputWidth, nil, nil).
		AddInputField("Club", "", inputWidth, nil, nil).
		AddButton("Add", func() {
			g.createPerson(form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "people")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
	//REVIEW: main or trips ? ^^
}

func (g *Gui) createPerson(form *tview.Form) {
	err := g.db.AddPerson(
		form.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Club").(*tview.InputField).GetText(),
	)
	if err != nil { // NOTE: Needs fixing
		g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
		return
	}

	g.closeAndSwitchPanel(`form`, `people`)
	g.app.QueueUpdateDraw(func() {
		g.peoplePanel().setEntries(g)
	})
}

//
// MODIFY FUNCS
func (g *Gui) modifyTripForm() {
	// First - what trip is selected?
	selectedTrip := g.selectedTrip()

	// Populate the form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Modify Trip ")
	form.SetTitleAlign(tview.AlignLeft)

	form.
		AddInputField("Date", selectedTrip.Date, inputWidth, nil, nil).
		AddInputField("Cave", selectedTrip.Cave, inputWidth, nil, nil).
		AddInputField("Names", selectedTrip.Names, inputWidth, nil, nil).
		AddInputField("Notes",  selectedTrip.Notes, inputWidth, nil, nil).
		AddButton("Apply", func() {
			g.modifyTrip(selectedTrip.ID, form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "trips")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) modifyTrip(id string, form *tview.Form) {
	err := g.db.ModifyTrip(
		id,
		form.GetFormItemByLabel("Date").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Cave").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Names").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Notes").(*tview.InputField).GetText(),
	)
	if err != nil {
		g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
		return
	}

	g.closeAndSwitchPanel(`form`, `trips`)
	g.app.QueueUpdateDraw(func() {
		g.tripsPanel().setEntries(g)
	})
}


func (g *Gui) modifyPersonForm() {
	// First - what trip is selected?
	selectedPerson := g.selectedPerson()

	// Populate the form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Modify Person ")
	form.SetTitleAlign(tview.AlignLeft)

	form.
		AddInputField("Name", selectedPerson.Name, inputWidth, nil, nil).
		AddInputField("Club", selectedPerson.Club, inputWidth, nil, nil).
		AddButton("Apply", func() {
			g.modifyPerson(selectedPerson.ID, form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "people")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) modifyPerson(id string, form *tview.Form) {
	err := g.db.ModifyPerson(
		id,
		form.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Club").(*tview.InputField).GetText(),
	)
	if err != nil {
		g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
		return
	}

	g.closeAndSwitchPanel(`form`, `people`)
	g.app.QueueUpdateDraw(func() {
		g.peoplePanel().setEntries(g)
	})
}


func (g *Gui) modifyLocationForm() {
	// First - what trip is selected?
	selectedLocation := g.selectedLocation()

	// Populate the form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Modify Location ")
	form.SetTitleAlign(tview.AlignLeft)

	form.
		AddInputField("Name", selectedLocation.Name, inputWidth, nil, nil).
		AddInputField("Region", selectedLocation.Region, inputWidth, nil, nil).
		AddInputField("Country", selectedLocation.Country, inputWidth, nil, nil).
		AddCheckbox("SRT", selectedLocation.SRT, nil).
		AddButton("Apply", func() {
			g.modifyLocation(selectedLocation.ID, form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "locations")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) modifyLocation(id string, form *tview.Form) {
	err := g.db.ModifyLocation(
		id,
		form.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Region").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Country").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("SRT").(*tview.Checkbox).IsChecked(),
	)
	if err != nil {
		g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
		return
	}

	g.closeAndSwitchPanel(`form`, `locations`)
	g.app.QueueUpdateDraw(func() {
		g.locationsPanel().setEntries(g)
	})
}


//
// DELETE FUNCS
func (g *Gui) deleteTrip() {
	selectedTrip := g.selectedTrip()

	message := fmt.Sprintf(
		"Are you sure you want to delete trip:\nDate: %s\nCave: %s",
		selectedTrip.Date, selectedTrip.Cave,
	)

	g.warning(message, `trips`, []string{`Yes`, `No`}, func() {
		if err := g.db.RemoveTrip(selectedTrip.ID); err != nil {
			g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
			return
		}
		g.tripsPanel().updateEntries(g)
	})
}

func (g *Gui) deleteLocation() {
	selectedLocation := g.selectedLocation()

	message := fmt.Sprintf(
		"Are you sure you want to delete location: %s",
		selectedLocation.Name,
	)

	g.warning(message, `locations`, []string{`Yes`, `No`}, func() {
		if err := g.db.RemoveLocation(selectedLocation.ID); err != nil {
			g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
			return
		}
		g.locationsPanel().updateEntries(g)
	})
}

func (g *Gui) deletePerson() {
	selectedPerson := g.selectedPerson()

	message := fmt.Sprintf(
		"Are you sure you want to delete person: %s",
		selectedPerson.Name,
	)

	g.warning(message, `people`, []string{`Yes`, `No`}, func() {
		if err := g.db.RemovePerson(selectedPerson.ID); err != nil {
			g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
			return
		}
		g.peoplePanel().updateEntries(g)
	})
}
