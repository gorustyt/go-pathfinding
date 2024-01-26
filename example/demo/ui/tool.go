package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/gorustyt/go-pathfinding/example/cmd/fyne_demo/tutorials"

	"fyne.io/fyne/v2/widget"
)

var toolData = map[string][]string{
	"": {
		"AStar",
		"Ida star",
		"Breadth First Search",
		"Dijkstra",
		"Jump Point Search",
		"Orthogonal Jump Point",
		"Trace",
	},
}
var (
	Options      = "Options"
	optionsGroup = []string{
		"Allow ",
		"Diagonal",
		"Bi-directional",
	}
)

func gewVisualizeRecursion() fyne.CanvasObject {
	return widget.NewCheckGroup([]string{"Visualize recursion"}, func(strings []string) {

	})
}

func getHeuristic() (res []fyne.CanvasObject) {
	return []fyne.CanvasObject{
		widget.NewLabel("Heuristic"),
		widget.NewRadioGroup([]string{
			"Manhattan ",
			"Euclidean ",
			"Octile",
			"Chebyshev",
		}, func(s string) {

		}),
	}
}

func getOptions3() (res []fyne.CanvasObject) {
	return []fyne.CanvasObject{
		widget.NewLabel("Options"),
		widget.NewCheckGroup([]string{
			"Allow Diagonal",
			"Bi-directional",
			"Don't Cross Corners",
		}, func(strings []string) {

		}),
	}
}

func getOptions2() (res []fyne.CanvasObject) {
	return []fyne.CanvasObject{
		widget.NewLabel("Options"),
		widget.NewCheckGroup([]string{
			"Allow Diagonal",
			"Bi-directional",
		}, func(strings []string) {

		}),
	}
}

func getAStarObjs() (res []fyne.CanvasObject) {
	res = append(getHeuristic(), getOptions3()...)
	res = append(res, container.NewGridWithColumns(2,
		widget.NewEntry(),
		widget.NewLabel("Weight")))
	return res
}

func getIDAStarObjs() (res []fyne.CanvasObject) {
	res = append(getHeuristic(), getOptions2()...)
	res = append(res, container.NewGridWithColumns(2,
		widget.NewEntry(),
		widget.NewLabel("Weight")))
	res = append(res, container.NewGridWithColumns(2,
		widget.NewEntry(),
		widget.NewLabel("Seconds limit")))
	res = append(res, gewVisualizeRecursion())
	return res
}

var toolObj = map[string]fyne.CanvasObject{
	"AStar":                container.NewHBox(getAStarObjs()...),
	"Ida star":             container.NewHBox(getIDAStarObjs()...),
	"Breadth First Search": container.NewHBox(getOptions3()...),
	"Best First Search":    container.NewHBox(append(getHeuristic(), getOptions3()...)...),
	"Dijkstra":             container.NewHBox(getOptions3()...),
	"Jump Point Search": container.NewHBox(
		append(getHeuristic(), gewVisualizeRecursion())...),
	"Orthogonal Jump Point": container.NewHBox(
		append(getHeuristic(), gewVisualizeRecursion())...),
	"Trace": container.NewHBox(append(getHeuristic(), getOptions2()...)...),
}

const preferenceCurrentTutorial = "currentTutorial"

func CreateToolTree(view func(w fyne.Window) fyne.CanvasObject) {
	a := fyne.CurrentApp()
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return tutorials.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := tutorials.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := tutorials.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := tutorials.Tutorials[uid]; ok {
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				view(t)
			}
		},
	}
	currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
	tree.Select(currentPref)
}
