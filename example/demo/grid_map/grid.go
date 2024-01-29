package grid_map

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
)

type grid struct {
	base        widget.BaseWidget
	g           *canvas.Rectangle
	i, j        float32
	scale       float32
	isHide      bool
	fillColor   color.NRGBA
	strokeColor color.NRGBA
	offset      fyne.Position
	m           *Map
	dragEndPos  fyne.Position

	popUpText *widget.Label
	popUp     *widget.PopUp
}

var (
	strokeColor = color.NRGBA{R: 255, G: 120, B: 0, A: 255}
	white       = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	gray        = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
)

func (g *grid) SetWalkAble(enable bool) {
	if !enable {
		g.fillColor = gray
		g.strokeColor = strokeColor
	} else {
		g.fillColor = white
		g.strokeColor = white
	}
}
func (g *grid) SetStart() {
	g.fillColor = color.NRGBA{G: 255, A: 255}
	g.strokeColor = strokeColor
}

func (g *grid) SetEnd() {
	g.fillColor = color.NRGBA{R: 255, A: 255}
	g.strokeColor = strokeColor
}

func (g *grid) SetPath() {
	g.fillColor = color.NRGBA{B: 255, A: 255}
	g.strokeColor = strokeColor
}

func (g *grid) Dragged(ev *fyne.DragEvent) {
	g.dragEndPos = ev.AbsolutePosition
}
func (g *grid) DragEnd() {

}

func (g *grid) MinSize() fyne.Size {
	return fyne.Size{Width: g.scale, Height: g.scale}
}

func (g *grid) MouseIn(ev *desktop.MouseEvent) {
	g.ShowPopUp(ev)
}

func (g *grid) ShowPopUp(ev *desktop.MouseEvent) {
	text := fmt.Sprintf("(%v,%v)", g.i, g.j)
	if g.popUp == nil {
		// 创建弹出框的内容
		g.popUpText = widget.NewLabel(text)
		// 创建并显示 Popover
		g.popUp = widget.NewPopUp(g.popUpText, g.m.win.Canvas())
	}
	g.popUpText.Text = text
	g.popUp.Move(fyne.Position{X: ev.AbsolutePosition.X, Y: ev.AbsolutePosition.Y})
	g.popUp.Show() // 将弹出框显示在按钮的中心位置
}

func (g *grid) MouseMoved(ev *desktop.MouseEvent) {}

func (g *grid) MouseOut() {
	g.popUp.Hide()
}

func (g *grid) Tapped(e *fyne.PointEvent) {
	g.m.OnGridWalkAble(g)
}

func (g *grid) TappedSecondary(e *fyne.PointEvent) {
	log.Println("I have been right tapped at", e, g.i, g.j)
}

func (g *grid) Move(position fyne.Position) {
	g.g.Move(position)
}

func (g *grid) Position() fyne.Position {
	return g.base.Position()
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
	g.g.FillColor = g.fillColor
	if g.isHide {
		g.g.Hide()
	} else {
		g.base.Move(fyne.Position{X: g.i * g.scale, Y: g.j * g.scale})
		g.g.Resize(fyne.NewSize(g.scale, g.scale))
		g.base.Resize(fyne.NewSize(g.scale, g.scale))
	}
	g.g.Refresh()

}

func (g *grid) CreateRenderer() fyne.WidgetRenderer {
	return newGridRenderer(g)
}

func newGrid(i, j int, scale float32) *grid {
	n := &grid{i: float32(i), j: float32(j), scale: scale, isHide: false, fillColor: white}
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

}

func (g gridRenderer) Layout(size fyne.Size) {
	g.Resize(size)
}

func (g gridRenderer) Objects() (res []fyne.CanvasObject) {
	return append(res, g.g)
}
