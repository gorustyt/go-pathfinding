package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	path_finding "github.com/gorustyt/go-pathfinding"
	"github.com/gorustyt/go-pathfinding/example/demo/grid_map"
)

var toolData = map[string][]string{
	"": {
		path_finding.DescAStar,
		path_finding.DescIdaStar,
		path_finding.DescBreadthFirstSearch,
		path_finding.DescBestFirstSearch,
		path_finding.DescDijkstra,
		path_finding.DescJumpPointSearch,
		path_finding.DescOrthogonalJumpPoint,
		path_finding.DescTrace,
	},
	path_finding.DescAStar:               {"none"},
	path_finding.DescIdaStar:             {"none"},
	path_finding.DescBreadthFirstSearch:  {"none"},
	path_finding.DescBestFirstSearch:     {"none"},
	path_finding.DescDijkstra:            {"none"},
	path_finding.DescJumpPointSearch:     {"none"},
	path_finding.DescOrthogonalJumpPoint: {"none"},
	path_finding.DescTrace:               {"none"},
}

func gewVisualizeRecursion() fyne.CanvasObject {
	g := widget.NewCheckGroup([]string{"Visualize recursion"}, func(strings []string) {

	})
	g.Selected = []string{"Visualize recursion"}
	return g
}

func getHeuristic(cfg *grid_map.Config) (res []fyne.CanvasObject) {
	g := widget.NewRadioGroup([]string{
		path_finding.DescHeuristicManhattan,
		path_finding.DescHeuristicEuclidean,
		path_finding.DescHeuristicOctile,
		path_finding.DescHeuristicChebyshev,
	}, func(s string) {
		cfg.Heuristic = path_finding.GetHeuristicByDesc(s)
	})
	g.Selected = path_finding.DescHeuristicManhattan
	return []fyne.CanvasObject{
		widget.NewLabel("Heuristic"),
		g,
	}
}

const (
	optionsAllowDiagonal    = "Allow Diagonal"
	optionsBiDirectional    = "Bi-directional"
	optionsDontCrossCorners = "Don't Cross Corners"
)

func parseOptions(s []string, cfg *grid_map.Config) {
	allowDiagonal := cfg.AllowDiagonal
	biDirectional := cfg.BiDirectional
	dontCrossCorners := cfg.DontCrossCorners
	for _, v := range s {
		switch v {
		case optionsAllowDiagonal:
			allowDiagonal = true
		case optionsBiDirectional:
			biDirectional = true
		case optionsDontCrossCorners:
			dontCrossCorners = true
		}
	}
	cfg.AllowDiagonal = allowDiagonal
	cfg.BiDirectional = biDirectional
	cfg.DontCrossCorners = dontCrossCorners
}
func getOptions3(cfg *grid_map.Config) (res []fyne.CanvasObject) {
	g := widget.NewCheckGroup([]string{
		optionsAllowDiagonal,
		optionsBiDirectional,
		optionsDontCrossCorners,
	}, func(s []string) {
		parseOptions(s, cfg)
	})
	if cfg.AllowDiagonal {
		g.Selected = append(g.Selected, optionsAllowDiagonal)
	}
	if cfg.BiDirectional {
		g.Selected = append(g.Selected, optionsBiDirectional)
	}
	if cfg.DontCrossCorners {
		g.Selected = append(g.Selected, optionsDontCrossCorners)
	}
	return []fyne.CanvasObject{
		widget.NewLabel("Options"),
		g,
	}
}

func getOptions2(cfg *grid_map.Config) (res []fyne.CanvasObject) {
	return []fyne.CanvasObject{
		widget.NewLabel("Options"),
		widget.NewCheckGroup([]string{
			optionsAllowDiagonal,
			optionsBiDirectional,
		}, func(s []string) {
			parseOptions(s, cfg)
		}),
	}
}

func getAStarObjs(cfg *grid_map.Config) (res []fyne.CanvasObject) {
	entry := widget.NewEntry()
	entry.Text = "1"
	entry.Bind(cfg.Weight)
	res = append(getHeuristic(cfg), getOptions3(cfg)...)
	res = append(res, container.NewGridWithColumns(2,
		entry,
		widget.NewLabel("Weight")))
	return res
}

func getIDAStarObjs(cfg *grid_map.Config) (res []fyne.CanvasObject) {
	entry1 := widget.NewEntry()
	entry1.Bind(cfg.Weight)
	entry1.Text = "1"
	entry2 := widget.NewEntry()
	entry2.Bind(cfg.SecondLimit)
	entry2.Text = "10"
	res = append(getHeuristic(cfg), getOptions2(cfg)...)
	res = append(res, container.NewGridWithColumns(2,
		entry1,
		widget.NewLabel("Weight")))
	res = append(res, container.NewGridWithColumns(2,
		entry2,
		widget.NewLabel("Seconds limit")))
	res = append(res, gewVisualizeRecursion())
	return res
}

func getToolObj(uid string, cfg *grid_map.Config) fyne.CanvasObject {
	switch uid {
	case path_finding.DescAStar:
		c := container.NewVBox(getAStarObjs(cfg)...)
		return c
	case path_finding.DescIdaStar:
		return container.NewVBox(getIDAStarObjs(cfg)...)
	case path_finding.DescBreadthFirstSearch:
		return container.NewVBox(getOptions3(cfg)...)
	case path_finding.DescBestFirstSearch:
		return container.NewVBox(append(getHeuristic(cfg), getOptions3(cfg)...)...)
	case path_finding.DescDijkstra:
		return container.NewVBox(getOptions3(cfg)...)
	case path_finding.DescJumpPointSearch:
		return container.NewVBox(
			append(getHeuristic(cfg), gewVisualizeRecursion())...)
	case path_finding.DescOrthogonalJumpPoint:
		return container.NewVBox(
			append(getHeuristic(cfg), gewVisualizeRecursion())...)
	case path_finding.DescTrace:
		return container.NewVBox(append(getHeuristic(cfg), getOptions2(cfg)...)...)
	}
	return nil
}

const preferenceCurrentTutorial = "currentTutorial"

func CreateTool(w fyne.Window, view func(), config *grid_map.Config) fyne.CanvasObject {
	bStart := &widget.Button{
		Text:       "start",
		Importance: widget.SuccessImportance,
		OnTapped: func() {
			config.OnStart()
		},
	}
	bPause := &widget.Button{
		Text:       "pause",
		Importance: widget.WarningImportance,
		OnTapped: func() {
			config.OnPause()
		},
	}
	bStop := &widget.Button{
		Text:       "stop",
		Importance: widget.DangerImportance,
		OnTapped: func() {
			config.OnStop()
		},
	}
	bClear := &widget.Button{
		Text:       "clear",
		Importance: widget.MediumImportance,
		OnTapped: func() {
			config.OnClear()
		},
	}
	return container.NewBorder(container.NewVBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(4, bStart, bPause, bStop, bClear),
		widget.NewSeparator(),
	), nil, nil, nil, createToolTree(w, view, config))
}

func createToolTree(w fyne.Window, view func(), cfg *grid_map.Config) fyne.CanvasObject {
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
				return getToolObj(path_finding.DescAStar, cfg)
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
			cfg.Uid = uid
			view()
		},
	}
	currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, path_finding.DescAStar)
	tree.Select(currentPref)
	return tree
}
