package main

import (
	"github.com/rivo/tview"
	"os/exec"
	//"fmt"
)


func main() {
	projectName := ""
	dc := ""
	app := tview.NewApplication()
	form := tview.NewForm().
		AddDropDown("Region", []string{"dc3-a", "dc4-a"}, 0, func(datacenter string, option int){
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
	modal := tview.NewModal().
			SetText(output).
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

func projectExists(app *tview.Application, project string, dc string) {
	cmd := exec.Command("openstack", "--os-cloud", dc, "project", "show", project)
	out, err := cmd.Output()
	if err != nil {
		displayError("The project doesn't exist!", app)
	} else {
		displayRessources(string(out), app)
	}
}
