package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	path_finding "github.com/gorustyt/go-pathfinding"
	"github.com/gorustyt/go-pathfinding/example/demo/grid_map"
	"strings"
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
	path_finding.DescAStar:               {"_" + path_finding.DescAStar},
	path_finding.DescIdaStar:             {"_" + path_finding.DescIdaStar},
	path_finding.DescBreadthFirstSearch:  {"_" + path_finding.DescBreadthFirstSearch},
	path_finding.DescBestFirstSearch:     {"_" + path_finding.DescBestFirstSearch},
	path_finding.DescDijkstra:            {"_" + path_finding.DescDijkstra},
	path_finding.DescJumpPointSearch:     {"_" + path_finding.DescJumpPointSearch},
	path_finding.DescOrthogonalJumpPoint: {"_" + path_finding.DescOrthogonalJumpPoint},
	path_finding.DescTrace:               {"_" + path_finding.DescTrace},
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
				return NewOptionsMenu(cfg)
			}
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			if branch {
				obj.(*widget.Label).SetText(uid)
			} else {
				obj.(*OptionsMenu).SetUid(strings.TrimPrefix(uid, "_"))
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
