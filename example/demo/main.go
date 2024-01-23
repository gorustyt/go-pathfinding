package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"time"
)

// TODO 图形化demo
func main() {
	a := app.New()
	w := a.NewWindow("Hello World")
	r := canvas.NewRectangle(color.Black)

	w.SetContent(r)
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			r.StrokeWidth = 1
			r.Move(fyne.Position{X: 100, Y: 100})
			r.Resize(fyne.NewSize(40, 40))
			r.Refresh()
		}
	}()
	//grid.NewGrid(30, 30, w)
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
