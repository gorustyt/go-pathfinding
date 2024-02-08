package grid_map

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	path_finding "github.com/gorustyt/go-pathfinding"
	"log"
	"strconv"
	"sync"
)

const (
	minOneGird = 40
	maxOneGird = 120
)

type mapStatus int

const (
	mapStatusNone mapStatus = iota
	mapStatusPathFinding
	mapStatusPause
	mapStatusDumpOrLoading
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

	start          *grid
	end            *grid
	paths          []*grid
	peekPaths      []*grid
	jumpPointsPath []*grid
	obs            []*grid

	Cfg        *Config
	trace      path_finding.DebugTrace
	traceLock  sync.Mutex
	state      mapStatus
	draggedEnd *grid
}

func (g *Map) OnClear() {
	for _, v := range g.obs {
		v.SetWalkAble(true)
	}

	g.obs = g.obs[:0]
	g.clearAllPaths()

	g.traceLock.Lock()
	if g.trace != nil {
		g.trace.Exit()
		g.trace = nil
	}

	g.traceLock.Unlock()
	g.state = mapStatusNone
	g.Refresh()
}

func (g *Map) OnGridMouseIn(v *grid) {
	if g.state != mapStatusNone || g.draggedEnd == nil {
		return
	}
	if g.draggedEnd == g.start {
		g.start.SetWalkAble(true)
		g.start = v
		g.draggedEnd = v
		g.Refresh()
	}
	if g.draggedEnd == g.end {
		g.end.SetWalkAble(true)
		g.end = v
		g.draggedEnd = v
		g.Refresh()
	}
}
func (g *Map) clearPaths() {
	for _, v := range g.paths {
		v.SetWalkAble(true)
	}
	g.paths = g.paths[:0]
}
func (g *Map) clearPeekPaths() {
	for _, v := range g.peekPaths {
		v.SetWalkAble(true)
	}
	g.peekPaths = g.peekPaths[:0]
}

func (g *Map) clearAllPaths() {
	g.clearPaths()
	g.clearPeekPaths()
	for _, v := range g.jumpPointsPath {
		v.SetWalkAble(true)
	}
	g.jumpPointsPath = g.jumpPointsPath[:0]
}

func (g *Map) OnStart() {
	if g.state == mapStatusPathFinding {
		return
	} else if g.state == mapStatusPause {
		g.state = mapStatusPathFinding
		g.traceLock.Lock()
		if g.trace != nil {
			g.trace.Start()
		}
		g.traceLock.Unlock()
		return
	}
	g.state = mapStatusPathFinding
	g.clearAllPaths()
	gw, gh := g.GetGridSize()
	t, c := g.Cfg.GetPathFindingConfig()

	go func() {
		m := path_finding.NewGridWithConfig(gw, gh, &c)
		g.traceLock.Lock()
		c.Trace = path_finding.NewDebugTrace()
		g.trace = c.Trace
		g.traceLock.Unlock()
		c.Trace.SetPathHandle(func(v *path_finding.TracePoint) {
			gg := g.grids[v.X][v.Y]
			if gg == g.start || gg == g.end {
				return
			}
			for _, tmp := range g.obs {
				if tmp == gg {
					return
				}
			}
			if v.IsJumpPoint {

				for i, tmp := range g.peekPaths {
					if tmp == gg {
						g.peekPaths = append(g.peekPaths[:i], g.peekPaths[i+1:]...)
					}
				}
				g.jumpPointsPath = append(g.jumpPointsPath, gg)
			} else {
				for _, tmp := range g.jumpPointsPath {
					if tmp == gg {
						return
					}
				}
				for _, tmp := range g.peekPaths {
					if tmp == gg {
						return
					}
				}
				g.peekPaths = append(g.peekPaths, gg)
			}
			g.Refresh()
		})
		for _, v := range g.obs {
			m.SetWalkableAt(int(v.i), int(v.j), false)
		}
		starX, startY := int(g.start.i), int(g.start.j)
		endX, endY := int(g.end.i), int(g.end.j)
		paths := m.PathFindingRoute(t)(starX, startY, endX, endY)
		g.traceLock.Lock()
		if g.trace != nil {
			g.trace.Exit()
			g.trace.Wait()
		}
		g.trace = nil
		g.traceLock.Unlock()
		g.renderPaths(paths)
		log.Printf("path-finding :%v exit....", t)
	}()
	g.Refresh()
}

func (g *Map) renderPaths(paths []*path_finding.PathPoint) {
	g.clearPaths()
	for _, v := range paths {
		gg := g.grids[v.X][v.Y]
		if gg == g.start || gg == g.end {
			continue
		}
		for i, tmp := range g.peekPaths {
			if tmp == gg {
				g.peekPaths = append(g.peekPaths[:i], g.peekPaths[i+1:]...)
			}
		}
		for i, tmp := range g.jumpPointsPath {
			if tmp == gg {
				g.jumpPointsPath = append(g.jumpPointsPath[:i], g.jumpPointsPath[i+1:]...)
			}
		}
		g.paths = append(g.paths, gg)
	}
	g.Refresh()
}
func (g *Map) OnPause() {
	if g.state != mapStatusPathFinding {
		return
	}
	g.state = mapStatusPause
	g.traceLock.Lock()
	if g.trace != nil {
		g.trace.Pause()
	}
	g.traceLock.Unlock()
}

func (g *Map) OnGridWalkAble(gg *grid) {
	if g.state != mapStatusNone {
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
	if g.state != mapStatusNone {
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

func (g *Map) SetStart(i, j int) {
	if g.start != nil {
		g.start.SetWalkAble(true)
	}
	g.start = g.grids[i][j]
}

func (g *Map) SetEnd(i, j int) {
	if g.end != nil {
		g.end.SetWalkAble(true)
	}
	g.end = g.grids[i][j]
}

func (g *Map) MinSize() fyne.Size {
	gw, gh := g.GetGridSize()
	return fyne.Size{Width: float32(gw) * minOneGird, Height: float32(gh) * minOneGird}
}

func (g *Map) Dragged(ev *fyne.DragEvent) {

}
func (g *Map) DragEnd() {
	g.Refresh()
}

func (g *Map) DoubleTapped(ev *fyne.PointEvent) {
	x, y := g.GetIndex(ev.AbsolutePosition)
	log.Println("I have been double tapped at", ev.AbsolutePosition.X, ev.AbsolutePosition.Y, x, y)

}

func (g *Map) Scrolled(ev *fyne.ScrollEvent) {
	if ev.Scrolled.DY > 0 {
		g.scale += g.scale * 0.35
	} else {
		g.scale -= g.scale * .5
	}

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
	g.w = int(size.Width)
	g.h = int(size.Height)
	w, h := g.GetGridSize()
	g.resize(w, h)
	g.Refresh()
}

func (g *Map) Size() fyne.Size {
	return fyne.Size{Width: float32(g.w), Height: float32(g.h)}
}

func (g *Map) Hide() {
	g.hide = true
	g.Refresh()
}
func (g *Map) GetScale() float32 {
	return g.scale * g.oneGird
}
func (g *Map) Visible() bool {
	return !g.hide
}

func (g *Map) initStartEnd() {
	g.SetStart(0, 0)
	w, h := g.GetGridSize()
	if w > 10 && h > 10 {
		g.SetEnd(10, 10)
		return
	}
	g.SetEnd(w/2, h/2)
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

	for _, v := range g.peekPaths {
		v.SetPeekPath()
	}
	for _, v := range g.jumpPointsPath {
		v.SetJumpPointPath()
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
		oneGird: 40}
	g.Cfg = NewConfig(g.OnStart, g.OnPause, g.OnClear)
	g.Cfg.Dump = g.Dump
	g.Cfg.Load = g.Load
	g.base.ExtendBaseWidget(g)
	w, h = g.GetGridSize()
	g.resize(w, h)
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
	OnClear func()
	Dump    func(closer fyne.URIWriteCloser)
	Load    func(fyne.URIReadCloser)
}

func NewConfig(OnStart func(),
	OnPause func(),
	OnClear func()) *Config {
	c := &Config{
		OnClear:           OnClear,
		OnPause:           OnPause,
		OnStart:           OnStart,
		Weight:            binding.NewString(),
		SecondLimit:       binding.NewString(),
		PathFindingConfig: path_finding.GetDefaultConfig(),
	}
	c.SecondLimit.Set("10")
	c.Weight.Set("1")
	return c
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
