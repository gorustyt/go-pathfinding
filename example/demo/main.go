package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/gorustyt/go-pathfinding/example/demo/grid_map"
	"github.com/gorustyt/go-pathfinding/example/demo/ui"
)

func main() {
	width, height := 900, 900
	a := app.NewWithID("path-finding")
	w := a.NewWindow("寻路算法")
	m := grid_map.NewMap(width, height)
	content := container.NewStack()
	setView := func() {
		content.Objects = []fyne.CanvasObject{m}
		content.Refresh()
	}
	gridMap := container.NewBorder(
		nil, nil, nil, nil, content)
	split := container.NewHSplit(ui.CreateTool(setView), gridMap)
	split.Offset = 0.2
	w.SetContent(split)
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	w.ShowAndRun()
}
