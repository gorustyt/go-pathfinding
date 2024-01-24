package grid_map

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"log"
)

type grid struct {
	base   widget.BaseWidget
	g      *canvas.Rectangle
	i, j   float32
	scale  float32
	isHide bool
}

func (g *grid) MinSize() fyne.Size {
	return fyne.Size{Width: 10, Height: 10}
}

func (g *grid) Tapped(e *fyne.PointEvent) {

}

func (g *grid) TappedSecondary(e *fyne.PointEvent) {
	log.Println("I have been right tapped at", e, g.i, g.j)
}

func (g *grid) DoubleTapped(e *fyne.PointEvent) {
	log.Println("I have been double tapped at", e, g.i, g.j)
}

func (g *grid) Move(position fyne.Position) {
	g.g.Move(position)
}

func (g *grid) Position() fyne.Position {
	return g.g.Position()
}

func (g *grid) Resize(size fyne.Size) {
	g.g.Resize(size)
}

func (g *grid) Size() fyne.Size {
	return fyne.Size{Width: g.scale, Height: g.scale}
}

func (g *grid) Hide() {
	g.isHide = true
	g.Refresh()
}

func (g *grid) Visible() bool {
	return !g.isHide
}

func (g *grid) Show() {
	g.isHide = false
	g.Refresh()
}

func (g *grid) Refresh() {
	if g.isHide {
		g.g.Hide()
	} else {
		g.g.Move(fyne.Position{X: g.i * g.scale, Y: g.j * g.scale})
		g.g.Show()
		g.g.Resize(fyne.NewSize(g.scale-1, g.scale-1))
	}
	g.g.Refresh()

}

func (g *grid) CreateRenderer() fyne.WidgetRenderer {
	return newGridRenderer(g)
}

func newGrid(i, j int, scale float32) *grid {
	n := &grid{i: float32(i), j: float32(j), scale: scale, isHide: true}
	n.base.ExtendBaseWidget(n)
	return n
}

func newGridRenderer(g *grid) *gridRenderer {
	return &gridRenderer{g}
}

type gridRenderer struct {
	*grid
}

func (g gridRenderer) Destroy() {
	fmt.Println("destroy")
}

func (g gridRenderer) Layout(size fyne.Size) {
	g.Resize(size)
}

func (g gridRenderer) Objects() (res []fyne.CanvasObject) {
	return append(res, g.g)
}
