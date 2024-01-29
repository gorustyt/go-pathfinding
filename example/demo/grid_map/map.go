package grid_map

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	path_finding "github.com/gorustyt/go-pathfinding"
	"image/color"
	"log"
	"strconv"
	"sync"
)

const (
	minOneGird = 40
	maxOneGird = 120
)

type Map struct {
	win     fyne.Window
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

	start *grid
	end   *grid
	paths []*grid
	obs   []*grid

	Cfg *Config

	isClear bool
}

func (g *Map) OnClear() {
	fmt.Println("on_clean")
	g.clear()
	g.Refresh()
}

func (g *Map) clear() {

	for _, v := range g.obs {
		v.SetWalkAble(true)
	}
	for _, v := range g.paths {
		v.SetWalkAble(true)
	}
	g.obs = g.obs[:0]
	g.paths = g.paths[:0]
	g.isClear = true
}

func (g *Map) OnStart() {
	fmt.Println("on_start")
	g.isClear = false
	g.paths = g.paths[:0]
	gw, gh := g.GetGridSize()
	t, c := g.Cfg.GetPathFindingConfig()
	m := path_finding.NewGridWithConfig(gw, gh, &c)
	for _, v := range g.obs {
		m.SetWalkableAt(int(v.i), int(v.j), false)
	}
	starX, startY := int(g.start.i), int(g.start.j)
	endX, endY := int(g.end.i), int(g.end.j)
	paths := m.PathFindingRoute(t)(starX, startY, endX, endY)
	for _, v := range paths {
		gg := g.grids[v.X][v.Y]
		if gg == g.start || gg == g.end {
			continue
		}
		g.paths = append(g.paths, gg)
	}
	g.Refresh()
}

func (g *Map) OnPause() {
	fmt.Println("on_pause")
}

func (g *Map) OnStop() {
	fmt.Println("on_stop")
}

func (g *Map) OnGridWalkAble(gg *grid) {
	if !g.isClear {
		return
	}
	if gg == g.start || gg == g.end {
		return
	}
	j := -1
	for i, v := range g.obs {
		if v == gg {
			j = i
			break
		}
	}
	if j != -1 {
		g.obs = append(g.obs[:j], g.obs[j+1:]...)
	} else {
		g.obs = append(g.obs, gg)
	}
	g.Refresh()
}

func (g *Map) OnGridChange(gg *grid, pos fyne.Position) {
	if !g.isClear {
		return
	}
	i, j := g.GetIndex(pos)
	changeGird := g.grids[i][j]
	if gg == g.start && changeGird != g.end {
		gg.SetWalkAble(true)
		g.start = changeGird
	} else if gg == g.end && changeGird != g.start {
		gg.SetWalkAble(true)
		g.end = changeGird
	}
	g.Refresh()
}

func (g *Map) SetStartPosition(pos fyne.Position) {
	i, j := g.GetIndex(pos)
	g.start = g.grids[i][j]
}

func (g *Map) SetEndPosition(pos fyne.Position) {
	i, j := g.GetIndex(pos)
	g.end = g.grids[i][j]
}

func (g *Map) MinSize() fyne.Size {
	gw, gh := g.GetGridSize()
	return fyne.Size{float32(gw) * minOneGird, float32(gh) * minOneGird}
}

func (g *Map) Dragged(ev *fyne.DragEvent) {
	//offsetX := float32(int(ev.Dragged.DX) % int(g.oneGird))
	//offsetY := float32(int(ev.Dragged.DY) % int(g.oneGird))
	//g.offset = g.offset.SubtractXY(offsetX, offsetY)
}
func (g *Map) DragEnd() {
	g.Refresh()
}

func (g *Map) DoubleTapped(ev *fyne.PointEvent) {
	x, y := g.GetIndex(ev.AbsolutePosition)
	log.Println("I have been double tapped at", ev.AbsolutePosition.X, ev.AbsolutePosition.Y, x, y)
	//g.ShowPopUp()
	//g.m.OnGridWalkAble(g)
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
	g.start.SetStart()
	g.end.SetEnd()
	for _, v := range g.obs {
		v.SetWalkAble(false)
	}
	for _, v := range g.paths {
		v.SetPath()
	}

	gw, gh := g.GetGridSize()
	for i := 0; i < gw; i++ {
		for j := 0; j < gh; j++ {
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
	i = int(pos.X) / int(g.oneGird*g.scale/2)
	j = int(pos.Y) / int(g.oneGird*g.scale/2)
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

func NewMap(w, h int, win fyne.Window) fyne.CanvasObject {
	g := &Map{w: w,
		win:     win,
		h:       h,
		scale:   1,
		isClear: true,
		oneGird: 40}
	g.Cfg = NewConfig(g.OnStart, g.OnPause, g.OnStop, g.OnClear)
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
				gg.m = g
				g.grids[i][j] = gg
			}
		}
	}
	g.start = g.grids[0][0]
	g.end = g.grids[10][10]
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

type Config struct {
	Weight      binding.String
	SecondLimit binding.String
	Uid         string
	*path_finding.PathFindingConfig

	OnStart func()
	OnPause func()
	OnStop  func()
	OnClear func()
}

func NewConfig(OnStart func(),
	OnPause func(),
	OnStop func(), OnClear func()) *Config {
	return &Config{
		OnStop:            OnStop,
		OnPause:           OnPause,
		OnStart:           OnStart,
		OnClear:           OnClear,
		Weight:            binding.NewString(),
		SecondLimit:       binding.NewString(),
		PathFindingConfig: path_finding.GetDefaultConfig(),
	}
}

func (cfg *Config) GetPathFindingConfig() (t path_finding.PathFindingType, c path_finding.PathFindingConfig) {
	c = *cfg.PathFindingConfig
	s, err := cfg.Weight.Get()
	if err != nil {
		panic(err)
	}
	if s != "" {
		cfg.PathFindingConfig.Weight, err = strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
	}

	s, err = cfg.SecondLimit.Get()
	if err != nil {
		panic(err)
	}
	if s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
		cfg.PathFindingConfig.IdAStarTimeLimit = v
	}
	t = path_finding.GetPathFindingType(cfg.Uid, &c)
	return
}
