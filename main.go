package main

import (
	"os/exec"

	"github.com/rivo/tview"
	//"fmt"
)

func main() {
	projectName := ""
	dc := ""
	app := tview.NewApplication()
	form := tview.NewForm().
		AddDropDown("Region", []string{"dc3-a", "dc4-a"}, 0, func(datacenter string, option int) {
			dc = datacenter
		}).
		AddInputField("Project", "", 20, nil, func(project string) {
			projectName = project
		}).
		AddButton("Save", func() {
			projectExists(app, projectName, dc)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func displayError(message string, app *tview.Application) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			}
		})
	if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
		panic(err)
	}
}

func displayRessources(output string, app *tview.Application) {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	ressources := []string{"server", "volume", "network", "subnet", "loadbalancer"}

	//newRessource := tview.NewButton("test").
	//	SetSelectedFunc(func() {
	//		//nothin
	//	})

	navigation := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	updateNav := func() {
		navigation.Clear()
		for _, ressource := range ressources {
			navigation.AddItem(ressource, "", 0, func() {
				app.Stop()
			})
		}

	}
	updateNav()

	grid := tview.NewGrid().
		SetRows(2, 0).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(navigation, 0, 0, 2, 1, 0, 0, false).
		AddItem(newPrimitive("main2"), 0, 1, 1, 1, 0, 0, false).
		AddItem(newPrimitive("main3"), 1, 1, 1, 1, 0, 0, false)

		// (newPrimitive("main3"), 1, 1, 1, 1, 0, 0, false)
		// 1 = selected Row
		// 1 = selected column
		// 1 = how many rows it uses
		// 1 = how many column it uses
		// 0 = minimum grid hight (characters)
		// 0 = minimum grid width (characters)

	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}

}

func projectExists(app *tview.Application, project string, dc string) {
	cmd := exec.Command("openstack", "--os-cloud", dc, "project", "show", project)
	out, err := cmd.Output()
	if err != nil {
		displayError("The project doesn't exist!", app)
	} else {
		displayRessources(string(out), app)
	}
}
