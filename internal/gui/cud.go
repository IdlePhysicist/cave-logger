package gui

import (
	"fmt"
	"time"
	"strings"

	tview "gitlab.com/tslocum/cview"
)

var inputWidth = 70

//
// CREATE FUNCS
func (g *Gui) createTripForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Add Trip ")
	form.SetTitleAlign(tview.AlignLeft)

	caveField := tview.NewInputField().
		SetLabel("Cave").
		SetFieldWidth(inputWidth)
	caveField.SetAutocompleteFunc(func(current string) (matches []string) {
		if len(current) == 0 {
			return
		}

		for _, location := range g.state.resources.locations {
			if strings.HasPrefix(strings.ToLower(location.Name), strings.ToLower(current)) {
				matches = append(matches, location.Name)
			}
		}

		return
	})

	form.
		AddInputField("Date", time.Now().Format(`2006-01-02`), inputWidth, nil, nil).
		AddFormItem(caveField).
		AddInputField("Names", "", inputWidth, nil, nil).
		AddInputField("Notes", "", inputWidth, nil, nil).
		AddButton("Add", func() {
			g.createTrip(form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "trips")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
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
	g.tripsPanel().updateEntries(g)
}

func (g *Gui) createLocationForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Add Location ")
	form.SetTitleAlign(tview.AlignLeft)

	regionField := tview.NewInputField().
		SetLabel("Region").
		SetFieldWidth(inputWidth)
	regionField.SetAutocompleteFunc(func(current string) (matches []string) {
		if len(current) == 0 {
			return
		}

		for _, region := range g.uniqueRegion(g.state.resources.locations) {
			if strings.HasPrefix(strings.ToLower(region), strings.ToLower(current)) {
				matches = append(matches, region)
			}
		}
		return
	})

	countryField := tview.NewInputField().
		SetLabel("Country").
		SetFieldWidth(inputWidth)
	countryField.SetAutocompleteFunc(func(current string) (matches []string) {
		if len(current) == 0 {
			return
		}

		for _, country := range g.uniqueCountry(g.state.resources.locations) {
			if strings.HasPrefix(strings.ToLower(country), strings.ToLower(current)) {
				matches = append(matches, country)
			}
		}
		return
	})

	form.
		AddInputField("Name", "", inputWidth, nil, nil).
		AddFormItem(regionField).
		AddFormItem(countryField).
		AddCheckBox("SRT", "", false, nil).
		AddButton("Add", func() {
			g.createLocation(form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "caves")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) createLocation(form *tview.Form) {
	err := g.db.AddLocation(
		form.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Region").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Country").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("SRT").(*tview.CheckBox).IsChecked(),
	)
	if err != nil { // NOTE: Needs fixing
		g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
		return
	}

	g.closeAndSwitchPanel(`form`, `caves`)
	g.cavesPanel().updateEntries(g)
}

func (g *Gui) createPersonForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Add Person ")
	form.SetTitleAlign(tview.AlignLeft)

	clubField := tview.NewInputField().
		SetLabel("Club").
		SetFieldWidth(inputWidth)
	clubField.SetAutocompleteFunc(func(current string) (matches []string) {
		if len(current) == 0 {
			return
		}

		for _, club := range g.uniqueClubs(g.state.resources.people) {
			if strings.HasPrefix(strings.ToLower(club), strings.ToLower(current)) {
				matches = append(matches, club)
			}
		}

		return
	})

	form.
		AddInputField("Name", "", inputWidth, nil, nil).
		AddFormItem(clubField).
		AddButton("Add", func() {
			g.createPerson(form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "cavers")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
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

	g.closeAndSwitchPanel(`form`, `cavers`)
	g.caversPanel().updateEntries(g)
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

	caveField := tview.NewInputField().
		SetLabel("Cave").
		SetFieldWidth(inputWidth).
		SetText(selectedTrip.Cave)
	caveField.SetAutocompleteFunc(func(current string) (matches []string) {
		if len(current) == 0 {
			return
		}

		for _, location := range g.state.resources.locations {
			if strings.HasPrefix(strings.ToLower(location.Name), strings.ToLower(current)) {
				matches = append(matches, location.Name)
			}
		}

		if len(matches) <=  1 ||  matches[0] == current {
			matches = nil
		}

		return
	})

	form.
		AddInputField("Date", selectedTrip.Date, inputWidth, nil, nil).
		AddFormItem(caveField).
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
	g.tripsPanel().updateEntries(g)
}


func (g *Gui) modifyPersonForm() {
	// First - what trip is selected?
	selectedPerson := g.selectedPerson()

	// Populate the form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Modify Person ")
	form.SetTitleAlign(tview.AlignLeft)

	clubField := tview.NewInputField().
		SetLabel("Club").
		SetFieldWidth(inputWidth).
		SetText(selectedPerson.Club)
	clubField.SetAutocompleteFunc(func(current string) (matches []string) {
		if len(current) == 0 {
			return
		}

		for _, club := range g.uniqueClubs(g.state.resources.people) {
			if strings.HasPrefix(strings.ToLower(club), strings.ToLower(current)) {
				matches = append(matches, club)
			}
		}

		return
	})

	form.
		AddInputField("Name", selectedPerson.Name, inputWidth, nil, nil).
		AddFormItem(clubField).
		AddButton("Apply", func() {
			g.modifyPerson(selectedPerson.ID, form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "cavers")
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

	g.closeAndSwitchPanel(`form`, `cavers`)
	g.caversPanel().updateEntries(g)
}


func (g *Gui) modifyLocationForm() {
	// First - what trip is selected?
	selectedLocation := g.selectedLocation()

	// Populate the form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Modify Location ")
	form.SetTitleAlign(tview.AlignLeft)

	regionField := tview.NewInputField().
		SetLabel("Region").
		SetFieldWidth(inputWidth).
		SetText(selectedLocation.Region)
	regionField.SetAutocompleteFunc(func(current string) (matches []string) {
		if len(current) == 0 {
			return
		}

		for _, region := range g.uniqueRegion(g.state.resources.locations) {
			if strings.HasPrefix(strings.ToLower(region), strings.ToLower(current)) {
				matches = append(matches, region)
			}
		}
		return
	})

	countryField := tview.NewInputField().
		SetLabel("Country").
		SetFieldWidth(inputWidth).
		SetText(selectedLocation.Country)
	countryField.SetAutocompleteFunc(func(current string) (matches []string) {
		if len(current) == 0 {
			return
		}

		for _, country := range g.uniqueCountry(g.state.resources.locations) {
			if strings.HasPrefix(strings.ToLower(country), strings.ToLower(current)) {
				matches = append(matches, country)
			}
		}
		return
	})


	form.
		AddInputField("Name", selectedLocation.Name, inputWidth, nil, nil).
		AddFormItem(regionField).
		AddFormItem(countryField).
		AddCheckBox("SRT", "", selectedLocation.SRT, nil).
		AddButton("Apply", func() {
			g.modifyLocation(selectedLocation.ID, form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "caves")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) modifyLocation(id string, form *tview.Form) {
	err := g.db.ModifyLocation(
		id,
		form.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Region").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("Country").(*tview.InputField).GetText(),
		form.GetFormItemByLabel("SRT").(*tview.CheckBox).IsChecked(),
	)
	if err != nil {
		g.warning(err.Error(), `form`, []string{`OK`}, func() {return})
		return
	}

	g.closeAndSwitchPanel(`form`, `caves`)
	g.cavesPanel().updateEntries(g)
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
		g.cavesPanel().updateEntries(g)
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
		g.caversPanel().updateEntries(g)
	})
}
