package gui

import (
	"context"
	"fmt"


	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/skanehira/docui/common"
)

var inputWidth = 70

func (g *Gui) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Rune() {
	case 'h':
		g.prevPanel()
	case 'l':
		g.nextPanel()
	case 'q':
		g.Stop()
	case '/':
		g.filter()
	}

	switch event.Key() {
	case tcell.KeyTab:
		g.nextPanel()
	case tcell.KeyBacktab:
		g.prevPanel()
	case tcell.KeyRight:
		g.nextPanel()
	case tcell.KeyLeft:
		g.prevPanel()
	}
}

func (g *Gui) filter() {
	currentPanel := g.state.panels.panel[g.state.panels.currentPanel]
	if currentPanel.name() == "tasks" {
		return
	}
	currentPanel.setFilterWord("")
	currentPanel.updateEntries(g)

	viewName := "filter"
	searchInput := tview.NewInputField().SetLabel("Word")
	searchInput.SetLabelWidth(6)
	searchInput.SetTitle("filter")
	searchInput.SetTitleAlign(tview.AlignLeft)
	searchInput.SetBorder(true)

	closeSearchInput := func() {
		g.closeAndSwitchPanel(viewName, g.state.panels.panel[g.state.panels.currentPanel].name())
	}

	searchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			closeSearchInput()
		}
	})

	searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			closeSearchInput()
		}
		return event
	})

	searchInput.SetChangedFunc(func(text string) {
		currentPanel.setFilterWord(text)
		currentPanel.updateEntries(g)
	})

	g.pages.AddAndSwitchToPage(viewName, g.modal(searchInput, 80, 3), true).ShowPage("main")
}

func (g *Gui) nextPanel() {
	idx := (g.state.panels.currentPanel + 1) % len(g.state.panels.panel)
	g.switchPanel(g.state.panels.panel[idx].name())
}

func (g *Gui) prevPanel() {
	g.state.panels.currentPanel--

	if g.state.panels.currentPanel < 0 {
		g.state.panels.currentPanel = len(g.state.panels.panel) - 1
	}

	idx := (g.state.panels.currentPanel) % len(g.state.panels.panel)
	g.switchPanel(g.state.panels.panel[idx].name())
}

func (g *Gui) addTripForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Add Log Entry")
	form.SetTitleAlign(tview.AlignLeft)

	form.AddInputField("Date", "", inputWidth, nil, nil).
		AddInputField("Cave", "", inputWidth, nil, nil).
		AddInputField("Names", "", inputWidth, nil, nil).
		AddInputField("Notes", "", inputWidth, nil, nil).
		AddButton("Add", func() {
			g.addTrip(form )
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "trips")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 29), true).ShowPage("main")
}

func (g *Gui) addTrip(form *tview.Form) {
	g.startTask("create container ", func(ctx context.Context) error {
		err := g.db.AddLog(
			form.GetFormItemByLabel("Date").(*tview.InputField).GetText(),
			form.GetFormItemByLabel("Cave").(*tview.InputField).GetText(),
			form.GetFormItemByLabel("Names").(*tview.InputField).GetText(),
			form.GetFormItemByLabel("Notes").(*tview.InputField).GetText(),
		)
		if err != nil {
			g.log.Errorf("cannot create entry %s", err)
			return err
		}


		g.closeAndSwitchPanel("form", "trips")
		g.app.QueueUpdateDraw(func() {
			g.tripPanel().setEntries(g)
		})

		return nil
	})
}

func (g *Gui) displayInspect(data, page string) {
	text := tview.NewTextView()
	text.SetTitle("Detail").SetTitleAlign(tview.AlignLeft)
	text.SetBorder(true)
	text.SetText(data)

	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			g.closeAndSwitchPanel("detail", page)
		}
		return event
	})

	g.pages.AddAndSwitchToPage("detail", text, true)
}

func (g *Gui) inspectTrip() {
	trip := g.selectedTrip()

	inspect, err := g.db.GetLogs(trip.ID)
	if err != nil {
		g.log.Errorf("cannot inspect container %s", err)
		return
	}

	g.displayInspect(common.StructToJSON(inspect), "trips")
}

/*func (g *Gui) inspectVolume() {
	volume := g.selectedVolume()

	inspect, err := docker.Client.InspectVolume(volume.Name)
	if err != nil {
		g.log.Errorf("cannot inspect volume %s", err)
		return
	}

	g.displayInspect(common.StructToJSON(inspect), "volumes")
}*/

func (g *Gui) removeTrip() {
	trip := g.selectedTrip()

	g.confirm("Do you want to remove the container?", "Done", "containers", func() {
		g.startTask(fmt.Sprintf("remove container %s", trip.Cave), func(ctx context.Context) error {
			if err := g.db.RemoveLog(trip.ID); err != nil {
				g.log.Errorf("cannot remove the container %s", err)
				return err
			}
			g.tripPanel().updateEntries(g)
			return nil
		})
	})
}

func (g *Gui) modifyTrip(form *tview.Form) {
	//trip := g.selectedTrip()

	
}

func (g *Gui) modifyTripForm() {

}

/*func (g *Gui) removeVolume() {
	volume := g.selectedVolume()

	g.confirm("Do you want to remove the volume?", "Done", "volumes", func() {
		g.startTask(fmt.Sprintf("remove volume %s", volume.Name), func(ctx context.Context) error {
			if err := docker.Client.RemoveVolume(volume.Name); err != nil {
				g.log.Errorf("cannot remove the volume %s", err)
				return err
			}
			g.volumePanel().updateEntries(g)
			return nil
		})
	})
}*/







/*

func (g *Gui) createVolumeForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Create volume")
	form.AddInputField("Name", "", inputWidth, nil, nil).
		AddInputField("Labels", "", inputWidth, nil, nil).
		AddInputField("Driver", "", inputWidth, nil, nil).
		AddInputField("Options", "", inputWidth, nil, nil).
		AddButton("Create", func() {
			g.createVolume(form)
		}).
		AddButton("Cancel", func() {
			g.closeAndSwitchPanel("form", "volumes")
		})

	g.pages.AddAndSwitchToPage("form", g.modal(form, 80, 13), true).ShowPage("main")
}

func (g *Gui) createVolume(form *tview.Form) {
	var data = make(map[string]string)
	inputLabels := []string{
		"Name",
		"Labels",
		"Driver",
		"Options",
	}

	for _, label := range inputLabels {
		data[label] = form.GetFormItemByLabel(label).(*tview.InputField).GetText()
	}

	g.startTask("create volume "+data["Name"], func(ctx context.Context) error {
		options := docker.Client.NewCreateVolumeOptions(data)

		if err := docker.Client.CreateVolume(options); err != nil {
			g.log.Errorf("cannot create volume %s", err)
			return err
		}

		g.closeAndSwitchPanel("form", "volumes")
		g.app.QueueUpdateDraw(func() {
			g.volumePanel().setEntries(g)
		})

		return nil
	})
}*/

