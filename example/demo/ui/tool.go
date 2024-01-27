package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	AStar               = "AStar"
	IdaStar             = "Ida star"
	BreadthFirstSearch  = "Breadth First Search"
	BestFirstSearch     = "Best First Search"
	Dijkstra            = "Dijkstra"
	JumpPointSearch     = "Jump Point Search"
	OrthogonalJumpPoint = "Orthogonal Jump Point"
	Trace               = "Trace"
)

var toolData = map[string][]string{
	"": {
		AStar,
		IdaStar,
		BreadthFirstSearch,
		BestFirstSearch,
		Dijkstra,
		JumpPointSearch,
		OrthogonalJumpPoint,
		Trace,
	},
	AStar:               {"none"},
	IdaStar:             {"none"},
	BreadthFirstSearch:  {"none"},
	BestFirstSearch:     {"none"},
	Dijkstra:            {"none"},
	JumpPointSearch:     {"none"},
	OrthogonalJumpPoint: {"none"},
	Trace:               {"none"},
}

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

func getToolObj(uid string) fyne.CanvasObject {
	switch uid {
	case AStar:
		c := container.NewVBox(getAStarObjs()...)
		return c
	case IdaStar:
		return container.NewVBox(getIDAStarObjs()...)
	case BreadthFirstSearch:
		return container.NewVBox(getOptions3()...)
	case BestFirstSearch:
		return container.NewVBox(append(getHeuristic(), getOptions3()...)...)
	case Dijkstra:
		return container.NewVBox(getOptions3()...)
	case JumpPointSearch:
		return container.NewVBox(
			append(getHeuristic(), gewVisualizeRecursion())...)
	case OrthogonalJumpPoint:
		return container.NewVBox(
			append(getHeuristic(), gewVisualizeRecursion())...)
	case Trace:
		return container.NewVBox(append(getHeuristic(), getOptions2()...)...)
	}
	return nil
}

const preferenceCurrentTutorial = "currentTutorial"

func CreateTool(view func()) fyne.CanvasObject {
	bStart := &widget.Button{
		Text:       "start",
		Importance: widget.SuccessImportance,
		OnTapped:   func() { fmt.Println("high importance button") },
	}
	bPause := &widget.Button{
		Text:       "pause",
		Importance: widget.WarningImportance,
		OnTapped:   func() { fmt.Println("tapped danger button") },
	}
	bStop := &widget.Button{
		Text:       "stop",
		Importance: widget.DangerImportance,
		OnTapped:   func() { fmt.Println("tapped warning button") },
	}
	return container.NewVBox(createToolTree(view), container.NewGridWithColumns(3, bStart, bPause, bStop))
}

func createToolTree(view func()) fyne.CanvasObject {
	a := fyne.CurrentApp()
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return toolData[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := toolData[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			if !branch {
				return getToolObj(AStar)
			}
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			if branch {
				obj.(*widget.Label).SetText(uid)
			} else {

			}

		},
		OnSelected: func(uid string) {
			view()
		},
	}
	currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, AStar)
	tree.Select(currentPref)
	return tree
}
