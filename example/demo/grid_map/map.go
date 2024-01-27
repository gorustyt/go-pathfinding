package grid_map

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"sync"
)

const (
	minOneGird = 40
	maxOneGird = 120
)

type Map struct {
	base    widget.BaseWidget
	w, h    int //地图宽和高
	hide    bool
	rows    [][]*canvas.Line //横线
	cols    [][]*canvas.Line //竖线
	grids   [][]*grid        //格子
	oneGird float32          //一个格子大小
	offset  fyne.Position    //位置偏移
	scale   float32          //缩放
	lock    sync.Mutex

	start fyne.Position   //起点
	end   fyne.Position   //终点
	paths []fyne.Position //寻路路径
	obs   []fyne.Position //障碍
}

func (g *Map) SetStart(x, y float64) {
	g.start = fyne.Position{X: float32(x), Y: float32(y)}
}

func (g *Map) SetEnd(x, y float64) {
	g.start = fyne.Position{X: float32(x), Y: float32(y)}
}

func (g *Map) SetPath() {

}

func (g *Map) SetNotWalk(x, y float64) {

}

func (g *Map) MinSize() fyne.Size {
	gw, gh := g.GetGridSize()
	return fyne.Size{float32(gw) * minOneGird, float32(gh) * minOneGird}
}

func (g *Map) Dragged(ev *fyne.DragEvent) {
	offsetX := float32(int(ev.Dragged.DX) % int(g.oneGird))
	offsetY := float32(int(ev.Dragged.DY) % int(g.oneGird))
	g.offset = g.offset.SubtractXY(offsetX, offsetY)
}
func (g *Map) DragEnd() {
	g.Refresh()
}

func (g *Map) Scrolled(ev *fyne.ScrollEvent) {
	if ev.Scrolled.DY > 0 {
		g.scale += g.scale * 0.35
	} else {
		g.scale -= g.scale * .5
	}

	fmt.Println("scale===", g.scale)
	if g.scale <= minOneGird/g.oneGird {
		g.scale = minOneGird / g.oneGird
	}
	if g.scale >= maxOneGird/g.oneGird {
		g.scale = maxOneGird / g.oneGird
	}
	g.Refresh()
}

func (g *Map) Move(position fyne.Position) {
	g.offset = fyne.NewPos(0-position.X, 0-position.Y)
	g.Refresh()
}

func (g *Map) Position() fyne.Position {
	return fyne.Position{}
}

func (g *Map) Resize(size fyne.Size) {

}

func (g *Map) Size() fyne.Size {
	return fyne.Size{Width: float32(g.w), Height: float32(g.h)}
}

func (g *Map) Hide() {
	g.hide = true
	g.Refresh()
}

func (g *Map) Visible() bool {
	return !g.hide
}

func (g *Map) drawMap() {
	i, j := g.GetIndex(g.start)
	g.grids[i][j].SetStart()
	i, j = g.GetIndex(g.end)
	g.grids[i][j].SetEnd()
	for _, v := range g.obs {
		i, j = g.GetIndex(v)
		g.grids[i][j].SetNotWalk()
	}
	for _, v := range g.paths {
		i, j = g.GetIndex(v)
		g.grids[i][j].SetPath()
	}

	gw, gh := g.GetGridSize()
	for i = 0; i < gw; i++ {
		for j = 0; j < gh; j++ {
			if i != gw {
				//画横线
				g.rows[i][j].Position1 = fyne.NewPos(float32(i)*g.oneGird*g.scale, float32(j)*g.oneGird*g.scale).Subtract(g.offset)
				g.rows[i][j].Position2 = fyne.NewPos(float32(i+1)*g.oneGird*g.scale, float32(j)*g.oneGird*g.scale).Subtract(g.offset)
				g.rows[i][j].Refresh()
			}
			if j != gh {
				//画竖线
				g.cols[i][j].Position1 = fyne.NewPos(float32(i)*g.oneGird*g.scale, float32(j)*g.oneGird*g.scale).Subtract(g.offset)
				g.cols[i][j].Position2 = fyne.NewPos(float32(i)*g.oneGird*g.scale, float32(j+1)*g.oneGird*g.scale).Subtract(g.offset)
				g.cols[i][j].Refresh()
			}
			if i < gw && j < gh {
				g.grids[i][j].scale = g.oneGird * g.scale
				g.grids[i][j].offset = g.offset
				//画格子
				g.grids[i][j].Refresh()
			}
		}
	}

}
func (g *Map) GetIndex(pos fyne.Position) (i, j int) {
	i = int(pos.X) / int(g.oneGird)
	j = int(pos.Y) / int(g.oneGird)
	return i, j
}
func (g *Map) Show() {
	g.hide = false
	g.Refresh()
}

func (g *Map) GetGridSize() (w int, h int) {
	return g.w / int(g.oneGird), g.h / int(g.oneGird)
}

func (g *Map) Refresh() {
	g.drawMap()
}

func NewMap(w, h int) fyne.CanvasObject {
	g := &Map{w: w,
		h:       h,
		scale:   1,
		start:   fyne.Position{X: 50, Y: 50},
		end:     fyne.Position{X: 400, Y: 600},
		oneGird: 40}
	g.base.ExtendBaseWidget(g)
	maxW, maxH := g.GetGridSize()
	g.rows = make([][]*canvas.Line, maxW+1)
	for i := range g.rows {
		g.rows[i] = make([]*canvas.Line, maxH+1)
	}

	g.cols = make([][]*canvas.Line, maxW+1)
	for i := range g.rows {
		g.cols[i] = make([]*canvas.Line, maxH+1)
	}

	g.grids = make([][]*grid, maxW)
	for i := range g.grids {
		g.grids[i] = make([]*grid, maxH)
	}

	for i := 0; i < maxW; i++ {
		for j := 0; j < maxH; j++ {
			line := canvas.NewLine(color.Black)
			line.StrokeWidth = 1
			g.rows[i][j] = line
			line1 := canvas.NewLine(color.Black)
			line1.StrokeWidth = 1
			g.cols[i][j] = line1

			r := canvas.NewRectangle(color.Black)
			r.StrokeWidth = 1
			if i < maxW && j < maxH {
				gg := newGrid(i, j, g.oneGird)
				gg.g = r
				g.grids[i][j] = gg
			}
		}
	}
	return g
}

func (g *Map) CreateRenderer() fyne.WidgetRenderer {
	return newMapRenderer(g)
}

func newMapRenderer(g *Map) *MapRenderer {
	return &MapRenderer{g}
}

type MapRenderer struct {
	*Map
}

func (g MapRenderer) Destroy() {
	fmt.Println("destroy")
}

func (g MapRenderer) Layout(size fyne.Size) {
	g.Resize(size)
}

func (g MapRenderer) Objects() (res []fyne.CanvasObject) {
	maxW, maxH := g.GetGridSize()
	for i := 0; i < maxW; i++ {
		for j := 0; j < maxH; j++ {
			res = append(res, g.rows[i][j])
			res = append(res, g.cols[i][j])
			res = append(res, g.grids[i][j])
		}
	}
	return res
}
