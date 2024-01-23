package grid

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

type Grid struct {
	base  widget.Icon
	w, h  int
	win   fyne.Window
	hide  bool
	gird1 [][]*canvas.Line
	gird2 [][]*canvas.Line
	scale float32
}

func (g Grid) CreateRenderer() fyne.WidgetRenderer {
	//TODO implement me
	panic("implement me")
}

func (g Grid) MinSize() fyne.Size {
	return fyne.Size{1, 1}
}

func (g Grid) Move(position fyne.Position) {

}

func (g Grid) Position() fyne.Position {
	return fyne.Position{0, 0}
}

func (g Grid) Resize(size fyne.Size) {

}

func (g Grid) Size() fyne.Size {
	return fyne.Size{Width: float32(g.w), Height: float32(g.h)}
}

func (g Grid) Hide() {
	g.hide = true
}

func (g Grid) Visible() bool {
	return !g.hide
}

func (g Grid) drawMap() {
	for i := 1; i < g.w; i++ { //画横线
		for j := 0; j < g.h; j++ {
			g.gird1[i][j].Position1 = fyne.Position{X: float32(i-1) * g.scale, Y: float32(j) * g.scale}
			g.gird1[i][j].Position2 = fyne.Position{X: float32(i) * g.scale, Y: float32(j) * g.scale}
			g.gird1[i][j].Refresh()
		}
	}
	for j := 1; j < g.h; j++ {
		for i := 0; i < g.w; i++ { //画竖线
			g.gird2[i][j].Position1 = fyne.Position{X: float32(i) * g.scale, Y: float32(j-1) * g.scale}
			g.gird2[i][j].Position2 = fyne.Position{X: float32(i) * g.scale, Y: float32(j) * g.scale}
			g.gird1[i][j].Refresh()
		}
	}

}

func (g Grid) Show() {
	g.hide = false
	g.Refresh()
}

func (g Grid) Refresh() {
	g.drawMap()
}

func NewGrid(w, h int, win fyne.Window) fyne.CanvasObject {
	g := &Grid{w: w, h: h, win: win, scale: 40}
	g.gird1 = make([][]*canvas.Line, w)
	for i := range g.gird1 {
		g.gird1[i] = make([]*canvas.Line, h)
	}

	g.gird2 = make([][]*canvas.Line, w)
	for i := range g.gird2 {
		g.gird2[i] = make([]*canvas.Line, h)
	}

	cn := container.NewWithoutLayout()
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			line := canvas.NewLine(color.Black)
			line.StrokeWidth = 1
			g.gird1[i][j] = line
			line1 := canvas.NewLine(color.Black)
			line1.StrokeWidth = 1
			g.gird2[i][j] = line1
			cn.Add(line)
			cn.Add(line1)
		}
	}
	g.base.ExtendBaseWidget(g)
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			g.Refresh()
		}
	}()
	win.SetContent(cn)
	return g
}
