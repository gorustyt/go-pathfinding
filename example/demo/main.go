package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/gorustyt/go-pathfinding/example/demo/grid_map"
)

func main() {
	width, height := 900, 900
	a := app.New()
	w := a.NewWindow("寻路算法")
	grid_map.NewMap(width, height, w)
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	w.ShowAndRun()
}
