package gui

import (
	"fmt"
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cview"
)

var inputWidth = 70

//
// CREATE FUNCS
func (g *Gui) createTripForm() {
	form := cview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Add Trip ")
	form.SetTitleAlign(cview.AlignLeft)

	caveField := cview.NewInputField()
	caveField.SetLabel("Cave")
	caveField.SetFieldWidth(inputWidth)
	caveField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		for _, location := range g.state.resources.locations {
			if strings.HasPrefix(strings.ToLower(location.Name), strings.ToLower(current)) {
				matches = append(matches, cview.NewListItem(location.Name))
			}
		}

		return
	})

	// Caver Field
	caverField := cview.NewInputField()
	caverField.SetLabel("Names")
	caverField.SetFieldWidth(inputWidth)
	caverField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		names := strings.Split(current, `,`)

		// Process preceeding names
		for _, n := range names {
			n = strings.TrimSpace(n)
		}

		var lastPart string
		if len(names) == 0 { // Then we're still on the first name
			lastPart = current
		} else {
			lastPart = names[len(names)-1] // To get the last name
		}

		lastPart = strings.TrimPrefix(lastPart, ` `) // We should still trim whitespace

		if len(lastPart) == 0 {
			return
		}

		for _, caver := range g.state.resources.people {
			if strings.HasPrefix(strings.ToLower(caver.Name), strings.ToLower(lastPart)) {
				matches = append(
					matches,
					cview.NewListItem(
						fmt.Sprintf(
							"%s%s", textFactory(names[0:len(names)-1], `, `), caver.Name,
						),
					),
				)
			}
		}

		return
	})

	form.AddInputField("Date", time.Now().Format(`2006-01-02`), inputWidth, nil, nil)
	form.AddFormItem(caveField)
	form.AddFormItem(caverField)
	form.AddInputField("Notes", "", inputWidth, nil, nil)
	form.AddButton("Add", func() {
		g.createTrip(form)
	})
	form.AddButton("Cancel", func() {
		g.closeAndSwitchPanel("form", "trips")
	})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true) //.ShowPage("main")
	//REVIEW: main or trips ? ^^
}

func (g *Gui) createTrip(form *cview.Form) {
	err := g.reg.AddTrip(
		form.GetFormItemByLabel("Date").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Cave").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Names").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Notes").(*cview.InputField).GetText(),
	)
	if err != nil { // NOTE: Needs fixing
		g.warning(err.Error(), `form`, []string{`OK`}, func() { return })
		return
	}

	g.closeAndSwitchPanel(`form`, `trips`)
	g.tripsPanel().updateEntries(g)
}

func (g *Gui) createLocationForm() {
	form := cview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Add Location ")
	form.SetTitleAlign(cview.AlignLeft)

	regionField := cview.NewInputField()
	regionField.SetLabel("Region")
	regionField.SetFieldWidth(inputWidth)
	regionField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		for _, region := range g.uniqueRegion(g.state.resources.locations) {
			if strings.HasPrefix(strings.ToLower(region), strings.ToLower(current)) {
				matches = append(matches, cview.NewListItem(region))
			}
		}
		return
	})

	countryField := cview.NewInputField()
	countryField.SetLabel("Country")
	countryField.SetFieldWidth(inputWidth)
	countryField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		for _, country := range g.uniqueCountry(g.state.resources.locations) {
			if strings.HasPrefix(strings.ToLower(country), strings.ToLower(current)) {
				matches = append(matches, cview.NewListItem(country))
			}
		}
		return
	})

	form.AddInputField("Name", "", inputWidth, nil, nil)
	form.AddFormItem(regionField)
	form.AddFormItem(countryField)
	form.AddCheckBox("SRT", "", false, nil)
	form.AddInputField("Notes", "", inputWidth, nil, nil)
	form.AddButton("Add", func() {
		g.createLocation(form)
	})
	form.AddButton("Cancel", func() {
		g.closeAndSwitchPanel("form", "caves")
	})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) createLocation(form *cview.Form) {
	err := g.reg.AddCave(
		form.GetFormItemByLabel("Name").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Region").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Country").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Notes").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("SRT").(*cview.CheckBox).IsChecked(),
	)
	if err != nil { // NOTE: Needs fixing
		g.warning(err.Error(), `form`, []string{`OK`}, func() { return })
		return
	}

	g.closeAndSwitchPanel(`form`, `caves`)
	g.cavesPanel().updateEntries(g)
}

func (g *Gui) createPersonForm() {
	form := cview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Add Person ")
	form.SetTitleAlign(cview.AlignLeft)

	clubField := cview.NewInputField()
	clubField.SetLabel("Club")
	clubField.SetFieldWidth(inputWidth)
	clubField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		for _, club := range g.uniqueClubs(g.state.resources.people) {
			if strings.HasPrefix(strings.ToLower(club), strings.ToLower(current)) {
				matches = append(matches, cview.NewListItem(club))
			}
		}

		return
	})

	form.AddInputField("Name", "", inputWidth, nil, nil)
	form.AddFormItem(clubField)
	form.AddInputField("Notes", "", inputWidth, nil, nil)
	form.AddButton("Add", func() {
		g.createPerson(form)
	})
	form.AddButton("Cancel", func() {
		g.closeAndSwitchPanel("form", "cavers")
	})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) createPerson(form *cview.Form) {
	err := g.reg.AddCaver(
		form.GetFormItemByLabel("Name").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Club").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Notes").(*cview.InputField).GetText(),
	)
	if err != nil { // NOTE: Needs fixing
		g.warning(err.Error(), `form`, []string{`OK`}, func() { return })
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
	form := cview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Modify Trip ")
	form.SetTitleAlign(cview.AlignLeft)

	caveField := cview.NewInputField()
	caveField.SetLabel("Cave")
	caveField.SetFieldWidth(inputWidth)
	caveField.SetText(selectedTrip.Cave)
	caveField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		for _, location := range g.state.resources.locations {
			if strings.HasPrefix(strings.ToLower(location.Name), strings.ToLower(current)) {
				matches = append(matches, cview.NewListItem(location.Name))
			}
		}

		if len(matches) <= 1 || matches[0].GetMainText() == current {
			matches = nil
		}

		return
	})

	// Caver Field
	caverField := cview.NewInputField()
	caverField.SetLabel("Names")
	caverField.SetFieldWidth(inputWidth)
	caverField.SetText(selectedTrip.Names)
	caverField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		names := strings.Split(current, `,`)

		// Process preceeding names
		for _, n := range names {
			n = strings.TrimSpace(n)
		}

		var lastPart string
		if len(names) == 0 { // Then we're still on the first name
			lastPart = current
		} else {
			lastPart = names[len(names)-1] // To get the last name
		}

		lastPart = strings.TrimPrefix(lastPart, ` `) // We should still trim whitespace

		if len(lastPart) == 0 {
			return
		}

		for _, caver := range g.state.resources.people {
			if strings.HasPrefix(strings.ToLower(caver.Name), strings.ToLower(lastPart)) {
				matches = append(
					matches,
					cview.NewListItem(
						fmt.Sprintf("%s%s", textFactory(names[0:len(names)-1], `, `), caver.Name),
					),
				)
			}
		}

		return
	})

	form.AddInputField("Date", selectedTrip.Date, inputWidth, nil, nil)
	form.AddFormItem(caveField)
	form.AddFormItem(caverField)
	form.AddInputField("Notes", selectedTrip.Notes, inputWidth, nil, nil)
	form.AddButton("Apply", func() {
		g.modifyTrip(selectedTrip.ID, form)
	})
	form.AddButton("Cancel", func() {
		g.closeAndSwitchPanel("form", "trips")
	})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) modifyTrip(id string, form *cview.Form) {
	err := g.reg.ModifyTrip(
		id,
		form.GetFormItemByLabel("Date").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Cave").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Names").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Notes").(*cview.InputField).GetText(),
	)
	if err != nil {
		g.warning(err.Error(), `form`, []string{`OK`}, func() { return })
		return
	}

	g.closeAndSwitchPanel(`form`, `trips`)
	g.tripsPanel().updateEntries(g)
}

func (g *Gui) modifyPersonForm() {
	// First - what trip is selected?
	selectedPerson := g.selectedPerson()

	// Populate the form
	form := cview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Modify Person ")
	form.SetTitleAlign(cview.AlignLeft)

	clubField := cview.NewInputField()
	clubField.SetLabel("Club")
	clubField.SetFieldWidth(inputWidth)
	clubField.SetText(selectedPerson.Club)
	clubField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		for _, club := range g.uniqueClubs(g.state.resources.people) {
			if strings.HasPrefix(strings.ToLower(club), strings.ToLower(current)) {
				matches = append(matches, cview.NewListItem(club))
			}
		}

		return
	})

	form.AddInputField("Name", selectedPerson.Name, inputWidth, nil, nil)
	form.AddFormItem(clubField)
	form.AddInputField("Notes", selectedPerson.Notes, inputWidth, nil, nil)
	form.AddButton("Apply", func() {
		g.modifyPerson(selectedPerson.ID, form)
	})
	form.AddButton("Cancel", func() {
		g.closeAndSwitchPanel("form", "cavers")
	})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) modifyPerson(id string, form *cview.Form) {
	err := g.reg.ModifyCaver(
		id,
		form.GetFormItemByLabel("Name").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Club").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Notes").(*cview.InputField).GetText(),
	)
	if err != nil {
		g.warning(err.Error(), `form`, []string{`OK`}, func() { return })
		return
	}

	g.closeAndSwitchPanel(`form`, `cavers`)
	g.caversPanel().updateEntries(g)
}

func (g *Gui) modifyLocationForm() {
	// First - what trip is selected?
	selectedLocation := g.selectedLocation()

	// Populate the form
	form := cview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Modify Location ")
	form.SetTitleAlign(cview.AlignLeft)

	regionField := cview.NewInputField()
	regionField.SetLabel("Region")
	regionField.SetFieldWidth(inputWidth)
	regionField.SetText(selectedLocation.Region)
	regionField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		for _, region := range g.uniqueRegion(g.state.resources.locations) {
			if strings.HasPrefix(strings.ToLower(region), strings.ToLower(current)) {
				matches = append(matches, cview.NewListItem(region))
			}
		}
		return
	})

	countryField := cview.NewInputField()
	countryField.SetLabel("Country")
	countryField.SetFieldWidth(inputWidth)
	countryField.SetText(selectedLocation.Country)
	countryField.SetAutocompleteFunc(func(current string) (matches []*cview.ListItem) {
		if len(current) == 0 {
			return
		}

		for _, country := range g.uniqueCountry(g.state.resources.locations) {
			if strings.HasPrefix(strings.ToLower(country), strings.ToLower(current)) {
				matches = append(matches, cview.NewListItem(country))
			}
		}
		return
	})

	form.AddInputField("Name", selectedLocation.Name, inputWidth, nil, nil)
	form.AddFormItem(regionField)
	form.AddFormItem(countryField)
	form.AddCheckBox("SRT", "", selectedLocation.SRT, nil)
	form.AddInputField("Notes", selectedLocation.Notes, inputWidth, nil, nil)
	form.AddButton("Apply", func() {
		g.modifyLocation(selectedLocation.ID, form)
	})
	form.AddButton("Cancel", func() {
		g.closeAndSwitchPanel("form", "caves")
	})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true)
}

func (g *Gui) modifyLocation(id string, form *cview.Form) {
	err := g.reg.ModifyCave(
		id,
		form.GetFormItemByLabel("Name").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Region").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Country").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("Notes").(*cview.InputField).GetText(),
		form.GetFormItemByLabel("SRT").(*cview.CheckBox).IsChecked(),
	)
	if err != nil {
		g.warning(err.Error(), `form`, []string{`OK`}, func() { return })
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
		if err := g.reg.RemoveTrip(selectedTrip.ID); err != nil {
			g.warning(err.Error(), `form`, []string{`OK`}, func() { return })
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
		if err := g.reg.RemoveCave(selectedLocation.ID); err != nil {
			g.warning(err.Error(), `form`, []string{`OK`}, func() { return })
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
		if err := g.reg.RemoveCaver(selectedPerson.ID); err != nil {
			g.warning(err.Error(), `form`, []string{`OK`}, func() { return })
			return
		}
		g.caversPanel().updateEntries(g)
	})
}

func textFactory(sl []string, sep string) string {
	s := ``
	numEl := len(sl)
	if numEl == 0 {
		return s
	}

	for i, el := range sl {
		el = strings.TrimSpace(el)
		s += el
		if i < numEl {
			s += sep
		}
	}
	return s
}
