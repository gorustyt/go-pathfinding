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
	a := app.New()
	w := a.NewWindow("寻路算法")
	grid_map.NewMap(width, height, w)
	content := container.NewStack()
	setView := func(view func(w fyne.Window) fyne.CanvasObject) {
		content.Objects = []fyne.CanvasObject{view(w)}
		content.Refresh()
	}
	tutorial := container.NewBorder(
		nil, nil, nil, nil, content)
	split := container.NewHSplit(ui.CreateToolTree(setView), tutorial)
	split.Offset = 0.2
	w.SetContent(split)
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	w.ShowAndRun()
}
