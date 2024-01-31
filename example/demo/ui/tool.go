package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	},
	path_finding.DescAStar:               {"_" + path_finding.DescAStar},
	path_finding.DescIdaStar:             {"_" + path_finding.DescIdaStar},
	path_finding.DescBreadthFirstSearch:  {"_" + path_finding.DescBreadthFirstSearch},
	path_finding.DescBestFirstSearch:     {"_" + path_finding.DescBestFirstSearch},
	path_finding.DescDijkstra:            {"_" + path_finding.DescDijkstra},
	path_finding.DescJumpPointSearch:     {"_" + path_finding.DescJumpPointSearch},
	path_finding.DescOrthogonalJumpPoint: {"_" + path_finding.DescOrthogonalJumpPoint},
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

const preferenceCurrentTutorial = "currentTutorial"

func CreateTool(w fyne.Window, view func(), config *grid_map.Config) fyne.CanvasObject {
	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
		}),
	)
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
	bClear := &widget.Button{
		Text:       "clear",
		Importance: widget.DangerImportance,
		OnTapped: func() {
			config.OnClear()
		},
	}
	label := widget.NewLabel("setting theme")
	label.Alignment = fyne.TextAlignCenter

	label1 := widget.NewLabel("start pathfinding")
	label1.Alignment = fyne.TextAlignCenter
	return container.NewBorder(container.NewVBox(
		label,
		widget.NewSeparator(),
		themes,
		widget.NewSeparator(),
		label1,
		container.NewGridWithColumns(3, bStart, bPause, bClear),
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
