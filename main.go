package main

import (
	"encoding/json"
	"os/exec"
	"sort"
	"strconv"

	"github.com/rivo/tview"
)

type Ressource struct {
	ID     string
	Name   string
	Status string
	Flavor string
	Size   int
}

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

func displayRessources(output string, app *tview.Application, project string, dc string) {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	grid := tview.NewGrid()
	ressources := []string{"server", "volume", "network", "subnet", "loadbalancer"}
	infoList := newPrimitive("List of ressources")
	infoShow := newPrimitive("Information regarding ressource")

	navigation := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	updateNav := func() {
		navigation.Clear()
		for _, ressource := range ressources {
			navigation.AddItem(ressource, "", 0, func() {
				listRessources(app, project, dc, ressource, grid, infoList, infoShow)
			})
		}

	}
	updateNav()

	grid.
		SetRows(10, 0).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(navigation, 0, 0, 2, 1, 0, 0, false).
		AddItem(infoList, 0, 1, 1, 1, 0, 0, false).
		AddItem(infoShow, 1, 1, 1, 1, 0, 0, false)

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
		displayRessources(string(out), app, project, dc)
	}
}

func listRessources(app *tview.Application, project string, dc string, ressource string, grid *tview.Grid, infoList tview.Primitive, infoShow tview.Primitive) {
	cmd := exec.Command("openstack", "--os-cloud", dc, ressource, "list", "-f", "json")
	out, err := cmd.Output()
	if err != nil {
		displayError("No ressource found!", app)
	} else {
		var parsedRessource []Ressource
		err := json.Unmarshal(out, &parsedRessource)
		if err != nil {
			displayError("Couldn't parse ressource!", app)
		} else {
			ressourceList := tview.NewList().
				ShowSecondaryText(false).
				SetHighlightFullLine(true)

			updateList := func() {
				ressourceList.Clear()
				for _, ressourceItem := range parsedRessource {
					displayName := ressourceItem.ID + " | " + ressourceItem.Name
					if ressourceItem.Flavor != "" {
						displayName = displayName + " | " + ressourceItem.Flavor
					}

					ressourceList.AddItem(displayName+"  ", "", 0, func() {
						showRessource(app, dc, ressource, ressourceItem.ID, grid, infoShow)
					})
				}
			}
			updateList()
			grid.RemoveItem(infoList)
			grid.RemoveItem(ressourceList)
			grid.AddItem(ressourceList, 0, 1, 1, 1, 0, 0, true)
		}
	}

}

func showRessource(app *tview.Application, dc string, ressource string, id string, grid *tview.Grid, infoShow tview.Primitive) {
	cmd := exec.Command("openstack", "--os-cloud", dc, ressource, "show", id, "-f", "json")
	out, err := cmd.Output()
	if err != nil {
		displayError("No ressource found!", app)
	} else {

		//Using for range loops on maps does not provide the same order at each loop
		//Using sorted indexer slice is therefore necessary
		data := make(map[string]any)
		err := json.Unmarshal(out, &data)

		if err != nil {
			panic(err)
			displayError("Couldn't parse ressource!", app)
		} else {

			var orderedData []string
			for k := range data {
				orderedData = append(orderedData, k)
			}
			sort.Strings(orderedData)
			ressourceShow := tview.NewTable().
				SetFixed(0, 2).
				SetBorders(true)

			r := 0
			c := 0

			for _, key := range orderedData {
				if c == 0 {
					ressourceShow.SetCell(r, c,
						tview.NewTableCell(key).
							SetAlign(tview.AlignCenter))

					if str, ok := data[key].(string); ok {
						ressourceShow.SetCell(r, c+1,
							tview.NewTableCell(str).
								SetAlign(tview.AlignCenter))
					} else {
						if num, ok := data[key].(int); ok {
							ressourceShow.SetCell(r, c,
								tview.NewTableCell(strconv.Itoa(num)).
									SetAlign(tview.AlignCenter))
						}
					}
				}
				r++

			}

			grid.RemoveItem(infoShow)
			grid.RemoveItem(ressourceShow)
			grid.AddItem(ressourceShow, 1, 1, 1, 1, 0, 0, false)
		}
	}

}
