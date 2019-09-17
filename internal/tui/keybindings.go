package tui

import (
	"context"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	//"github.com/idlephysicist/cave-logger/internal/pkg/keeper"
)

var inputWidth = 70

func (t *Tui) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Rune() {
	case 'h':
		t.prevPanel()
	case 'l':
		t.nextPanel()
	case 'q':
		t.Stop()
	case '/':
		t.filter()
	}

	switch event.Key() {
	case tcell.KeyTab:
		t.nextPanel()
	case tcell.KeyBacktab:
		t.prevPanel()
	case tcell.KeyRight:
		t.nextPanel()
	case tcell.KeyLeft:
		t.prevPanel()
	}
}

func (t *Tui) filter() {
	currentPanel := t.state.panels.panel[t.state.panels.currentPanel]
	if currentPanel.name() == "tasks" {
		return
	}
	currentPanel.setFilterWord("")
	currentPanel.updateEntries(t)

	viewName := "filter"
	searchInput := tview.NewInputField().SetLabel("Word")
	searchInput.SetLabelWidth(6)
	searchInput.SetTitle("filter")
	searchInput.SetTitleAlign(tview.AlignLeft)
	searchInput.SetBorder(true)

	closeSearchInput := func() {
		t.closeAndSwitchPanel(viewName, t.state.panels.panel[t.state.panels.currentPanel].name())
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
		currentPanel.updateEntries(t)
	})

	t.pages.AddAndSwitchToPage(viewName, t.modal(searchInput, 80, 3), true).ShowPage("main")
}

func (t *Tui) nextPanel() {
	idx := (t.state.panels.currentPanel + 1) % len(t.state.panels.panel)
	t.switchPanel(t.state.panels.panel[idx].name())
}

func (t *Tui) prevPanel() {
	t.state.panels.currentPanel--

	if t.state.panels.currentPanel < 0 {
		t.state.panels.currentPanel = len(t.state.panels.panel) - 1
	}

	idx := (t.state.panels.currentPanel) % len(t.state.panels.panel)
	t.switchPanel(t.state.panels.panel[idx].name())
}

func (t *Tui) addEntryForm() {
	selectedEntry := t.selectedEntry()
	if selectedEntry == nil {
		//common.Logger.Error("please input image")
		return
	}

	//image := fmt.Sprintf("%s:%s", selectedEntry.Repo, selectedEntry.Tag)

	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Add Entry")
	form.SetTitleAlign(tview.AlignLeft)

	form.AddInputField("Date", "", inputWidth, nil, nil).
		AddInputField("Cave", "", inputWidth, nil, nil).
		AddInputField("Names", "", inputWidth, nil, nil).
		AddInputField("Notes", "", inputWidth, nil, nil).
		AddButton("Add", func() {
			t.addEntry(form, ``)
		}).
		AddButton("Cancel", func() {
			t.closeAndSwitchPanel("form", "images")
		})

	t.pages.AddAndSwitchToPage("form", t.modal(form, 80, 29), true).ShowPage("main")
}

func (t *Tui) addEntry(form *tview.Form, image string) {
	inputLabels := []string{
		"Date",
		"Cave",
		"Name",
		"Notes",
	}

	var params []string
	for _, data := range inputLabels {
		params = append(
			params, 
			form.GetFormItemByLabel(data).(*tview.InputField).GetText(),
		)
	}

	err := t.db.AddLog(params)
	if err != nil {
		//common.Logger.Errorf("cannot create container %s", err)
		return err
	}

	t.closeAndSwitchPanel("form", "images")
	t.app.QueueUpdateDraw(func() {
		t.containerPanel().setEntries(t)
	})

	return nil
}


func (t *Tui) displayInspect(data, page string) {
	text := tview.NewTextView()
	text.SetTitle("Detail").SetTitleAlign(tview.AlignLeft)
	text.SetBorder(true)
	text.SetText(data)

	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			t.closeAndSwitchPanel("detail", page)
		}
		return event
	})

	t.pages.AddAndSwitchToPage("detail", text, true)
}

func (t *Tui) inspectEntry() {
	entry := t.selectedEntry()

	inspect, err := t.db.QueryLogs(entry.ID)
	if err != nil {
		//common.Logger.Errorf("cannot inspect image %s", err)
		return
	}

	t.displayInspect(common.StructToJSON(inspect), "images")
}

func (t *Tui) removeEntry() {
	image := t.selectedEntry()

	t.confirm("Do you want to remove the image?", "Done", "images", func() {
		t.startTask(fmt.Sprintf("remove image %s:%s", image.Repo, image.Tag), func(ctx context.Context) error {
			if err := docker.Client.RemoveImage(image.ID); err != nil {
				common.Logger.Errorf("cannot remove the image %s", err)
				return err
			}
			t.imagePanel().updateEntries(g)
			return nil
		})
	})
}

func (t *Tui) removeContainer() {
	container := t.selectedContainer()

	t.confirm("Do you want to remove the container?", "Done", "containers", func() {
		t.startTask(fmt.Sprintf("remove container %s", container.Name), func(ctx context.Context) error {
			if err := docker.Client.RemoveContainer(container.ID); err != nil {
				common.Logger.Errorf("cannot remove the container %s", err)
				return err
			}
			t.containerPanel().updateEntries(g)
			return nil
		})
	})
}







func (t *Tui) loadImageForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Load image")
	form.AddInputField("Path", "", inputWidth, nil, nil).
		AddButton("Load", func() {
			path := form.GetFormItemByLabel("Path").(*tview.InputField).GetText()
			t.loadImage(path)
		}).
		AddButton("Cancel", func() {
			t.closeAndSwitchPanel("form", "images")
		})

	t.pages.AddAndSwitchToPage("form", t.modal(form, 80, 7), true).ShowPage("main")
}

func (t *Tui) loadImage(path string) {
	t.startTask("load image "+filepath.Base(path), func(ctx context.Context) error {
		t.closeAndSwitchPanel("form", "images")
		if err := docker.Client.LoadImage(path); err != nil {
			common.Logger.Errorf("cannot load image %s", err)
			return err
		}

		t.imagePanel().updateEntries(g)
		return nil
	})
}

func (t *Tui) importImageForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Import image")
	form.AddInputField("Repository", "", inputWidth, nil, nil).
		AddInputField("Tag", "", inputWidth, nil, nil).
		AddInputField("Path", "", inputWidth, nil, nil).
		AddButton("Load", func() {
			repository := form.GetFormItemByLabel("Repository").(*tview.InputField).GetText()
			tag := form.GetFormItemByLabel("Tag").(*tview.InputField).GetText()
			path := form.GetFormItemByLabel("Path").(*tview.InputField).GetText()
			t.importImage(path, repository, tag)
		}).
		AddButton("Cancel", func() {
			t.closeAndSwitchPanel("form", "images")
		})

	t.pages.AddAndSwitchToPage("form", t.modal(form, 80, 11), true).ShowPage("main")
}

func (t *Tui) importImage(file, repo, tag string) {
	t.startTask("import image "+file, func(ctx context.Context) error {
		t.closeAndSwitchPanel("form", "images")

		if err := docker.Client.ImportImage(repo, tag, file); err != nil {
			common.Logger.Errorf("cannot load image %s", err)
			return err
		}

		t.imagePanel().updateEntries(g)
		return nil
	})
}

func (t *Tui) saveImageForm() {
	image := t.selectedEntry()
	imageName := fmt.Sprintf("%s:%s", image.Repo, image.Tag)

	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Save image")
	form.AddInputField("Path", "", inputWidth, nil, nil).
		AddInputField("Image", imageName, inputWidth, nil, nil).
		AddButton("Save", func() {
			image := form.GetFormItemByLabel("Image").(*tview.InputField).GetText()
			path := form.GetFormItemByLabel("Path").(*tview.InputField).GetText()
			t.saveImage(image, path)
		}).
		AddButton("Cancel", func() {
			t.closeAndSwitchPanel("form", "images")
		})

	t.pages.AddAndSwitchToPage("form", t.modal(form, 80, 9), true).ShowPage("main")

}

func (t *Tui) saveImage(image, path string) {
	t.startTask("save image "+image, func(ctx context.Context) error {
		t.closeAndSwitchPanel("form", "images")

		if err := docker.Client.SaveImage([]string{image}, path); err != nil {
			common.Logger.Errorf("cannot save image %s", err)
			return err
		}
		return nil
	})

}



func (t *Tui) tailContainerLog() {
	container := t.selectedContainer()
	if container == nil {
		common.Logger.Errorf("cannot start tail container: selected container is null")
		return
	}

	if !t.app.Suspend(func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		errCh := make(chan error)

		go func() {
			reader, err := docker.Client.ContainerLogStream(container.ID)
			if err != nil {
				common.Logger.Error(err)
				errCh <- err
			}
			defer reader.Close()

			_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
			if err != nil {
				common.Logger.Error(err)
				errCh <- err
			}
			return
		}()

		select {
		case err := <-errCh:
			common.Logger.Error(err)
			return
		case <-sigint:
			return
		}
	}) {
		common.Logger.Error("cannot suspend tview")
	}
}