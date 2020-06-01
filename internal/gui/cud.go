package gui

import (
  "fmt"
  "time"
  "strings"

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

  caveField := tview.NewInputField().
    SetLabel("Cave").
		SetAutocompleteMultipleEntries(false).
    SetFieldWidth(inputWidth)
  caveField.SetAutocompleteFunc(func(current string) (matches []interface{}) {
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

	// Caver Field
	caverField := tview.NewInputField().
		SetLabel("Names").
		SetAutocompleteMultipleEntries(true).
		SetFieldWidth(inputWidth)
	caverField.SetAutocompleteFunc(func(current string) (matches []interface{}) {
		if len(current) == 0 {
			return []interface{}{[2]string{``,``}}
		}

		names := strings.Split(current, `,`)

		// Process preceeding names
		for _, n := range names {
			n = strings.TrimPrefix(n, ` `)
		}

		var lastPart string
		if len(names) == 0 {	// Then we're still on the first name
			lastPart = current
		} else {
			lastPart = names[len(names)-1]	// To get the last name
		}

		lastPart = strings.TrimPrefix(lastPart, ` `)	// We should still trim whitespace

		if len(lastPart) == 0 {
			return []interface{}{[2]string{``,``}}
		}

		for _, caver := range g.state.resources.people {
			if strings.HasPrefix(strings.ToLower(caver.Name), strings.ToLower(lastPart)) {
				matches = append(
					matches,
					[2]string{
						caver.Name,
						textFactory(names[0:len(names)-1], `, `),
					},
				)
			}
		}

		return
	})


  form.
    AddInputField("Date", time.Now().Format(`2006-01-02`), inputWidth, nil, nil).
    AddFormItem(caveField).
    AddFormItem(caverField).
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

  regionField := tview.NewInputField().
    SetLabel("Region").
		SetAutocompleteMultipleEntries(false).
    SetFieldWidth(inputWidth)
  regionField.SetAutocompleteFunc(func(current string) (matches []interface{}) {
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
		SetAutocompleteMultipleEntries(false).
    SetFieldWidth(inputWidth)
  countryField.SetAutocompleteFunc(func(current string) (matches []interface{}) {
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

  clubField := tview.NewInputField().
    SetLabel("Club").
		SetAutocompleteMultipleEntries(false).
    SetFieldWidth(inputWidth)
  clubField.SetAutocompleteFunc(func(current string) (matches []interface{}) {
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

  caveField := tview.NewInputField().
    SetLabel("Cave").
    SetFieldWidth(inputWidth).
		SetAutocompleteMultipleEntries(false).
    SetText(selectedTrip.Cave)
  caveField.SetAutocompleteFunc(func(current string) (matches []interface{}) {
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

  clubField := tview.NewInputField().
    SetLabel("Club").
    SetFieldWidth(inputWidth).
		SetAutocompleteMultipleEntries(false).
    SetText(selectedPerson.Club)
  clubField.SetAutocompleteFunc(func(current string) (matches []interface{}) {
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
    g.peoplePanel().updateEntries(g)//setEntries(g)
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

  regionField := tview.NewInputField().
    SetLabel("Region").
    SetFieldWidth(inputWidth).
		SetAutocompleteMultipleEntries(false).
    SetText(selectedLocation.Region)
  regionField.SetAutocompleteFunc(func(current string) (matches []interface{}) {
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
		SetAutocompleteMultipleEntries(false).
    SetText(selectedLocation.Country)
  countryField.SetAutocompleteFunc(func(current string) (matches []interface{}) {
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


func textFactory(sl []string, sep string) string {
	s := ``
	numEl := len(sl)
	if numEl == 0 {
		return s
	}

	for i, el := range sl {
		s += el
		if i < numEl {
			s += sep
		}
	}
	return s
}
