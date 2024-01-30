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
	widget.BaseWidget
	g           *canvas.Rectangle
	i, j        float32
	fillColor   color.NRGBA
	strokeColor color.NRGBA
	m           *Map
	dragEndPos  fyne.Position

	popUpText *widget.Label
	popUp     *widget.PopUp
}

var (
	strokeColor = color.NRGBA{R: 255, G: 120, B: 0, A: 255}
	white       = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	gray        = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
	blue        = color.NRGBA{B: 255, A: 255}
	yellow      = color.NRGBA{R: 255, G: 255, A: 150}
	green       = color.NRGBA{G: 255, A: 255}
	red         = color.NRGBA{R: 255, A: 255}
	tarquoise   = color.NRGBA{R: 64, G: 224, B: 208, A: 255}
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
func (g *grid) SetJumpPointPath() {
	g.fillColor = tarquoise
}
func (g *grid) SetPeekPath() {
	g.fillColor = yellow
}
func (g *grid) SetStart() {
	g.fillColor = green
	g.strokeColor = strokeColor
}

func (g *grid) SetEnd() {
	g.fillColor = red
	g.strokeColor = strokeColor
}

func (g *grid) SetPath() {
	g.fillColor = blue
	g.strokeColor = strokeColor
}

func (g *grid) Dragged(ev *fyne.DragEvent) {
	if g.m.draggedEnd == nil {
		g.m.draggedEnd = g
	}
}

func (g *grid) DragEnd() {
	g.m.draggedEnd = nil
}

func (g *grid) MinSize() fyne.Size {
	return fyne.Size{Width: g.m.GetScale(), Height: g.m.GetScale()}
}

func (g *grid) MouseIn(ev *desktop.MouseEvent) {
	g.m.OnGridMouseIn(g)
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
	//g.popUp.Hide()
}

func (g *grid) Tapped(e *fyne.PointEvent) {
	g.m.OnGridWalkAble(g)
}

func (g *grid) TappedSecondary(e *fyne.PointEvent) {
	log.Println("I have been right tapped at", e, g.i, g.j)
}

func (g *grid) Resize(size fyne.Size) {
	g.g.Resize(size)
}

func (g *grid) Size() fyne.Size {
	return fyne.Size{Width: g.m.GetScale(), Height: g.m.GetScale()}
}

func (g *grid) Refresh() {
	g.g.FillColor = g.fillColor
	g.Move(fyne.Position{X: g.i * g.m.GetScale(), Y: g.j * g.m.GetScale()})
	g.g.Resize(fyne.NewSize(g.m.GetScale(), g.m.GetScale()))
	g.Resize(fyne.NewSize(g.m.GetScale(), g.m.GetScale()))
	g.g.Refresh()

}

func (g *grid) CreateRenderer() fyne.WidgetRenderer {
	return newGridRenderer(g)
}

func newGrid(i, j int, scale float32) *grid {
	n := &grid{i: float32(i), j: float32(j), fillColor: white}
	n.ExtendBaseWidget(n)
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
