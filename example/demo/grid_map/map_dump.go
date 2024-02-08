package grid_map

import (
	"bufio"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"image/color"
	"log/slog"
)

const (
	mapStrWalkable = '0'
	mapStrObs      = '1'
	mapStrStart    = '2'
	mapStrEnd      = '3'
)

func (g *Map) resize(w, h int, initStartEnds ...bool) {
	oldW := 0
	if len(g.grids) > 0 {
		oldW = len(g.grids)
	}
	if oldW < w {
		newGrids := make([][]*grid, w)
		copy(newGrids, g.grids)
		g.grids = newGrids

		newRows := make([][]*canvas.Line, w)
		copy(newRows, g.rows)
		g.rows = newRows

		newCols := make([][]*canvas.Line, w)
		copy(newCols, g.cols)
		g.cols = newCols
	} else {
		g.grids = g.grids[:w]
		g.rows = g.rows[:w]
		g.cols = g.cols[:w]
	}

	for i := range g.grids {
		if len(g.grids[i]) < h {
			newGirds := make([]*grid, h)
			copy(newGirds, g.grids[i])

			newCols := make([]*canvas.Line, h)
			copy(newCols, g.cols[i])

			newRows := make([]*canvas.Line, h)
			copy(newRows, g.rows[i])

			for j := len(g.grids[i]); j < h; j++ {
				line := canvas.NewLine(color.Black)
				line.StrokeWidth = 1
				newRows[j] = line
				line1 := canvas.NewLine(color.Black)
				line1.StrokeWidth = 1
				newCols[j] = line1

				r := canvas.NewRectangle(color.Black)
				r.StrokeWidth = 1

				gg := newGrid(i, j)
				gg.g = r
				gg.m = g
				newGirds[j] = gg
			}
			g.grids[i] = newGirds
			g.rows[i] = newRows
			g.cols[i] = newCols
		} else {
			g.grids[i] = g.grids[i][:h]
			g.cols[i] = g.cols[i][:h]
			g.rows[i] = g.rows[i][:h]
		}
	}
	forceInitStartEnd := false
	if g.start != nil && g.end != nil && (g.start.i >= float32(w) || g.start.j >= float32(h) || g.end.i >= float32(w) || g.end.j >= float32(h)) {
		forceInitStartEnd = true
	}
	if forceInitStartEnd || len(initStartEnds) == 0 || (len(initStartEnds) > 0 && initStartEnds[0]) {
		g.initStartEnd()
	}
}

func (g *Map) Dump(f fyne.URIWriteCloser) {
	oldState := g.state
	if g.state == mapStatusDumpOrLoading {
		return
	}
	if g.state == mapStatusPathFinding {
		g.OnPause()
	}
	g.state = mapStatusDumpOrLoading
	gw, gh := g.GetGridSize()
	ms := make([][]byte, gw)
	for i := 0; i < gw; i++ {
		for j := 0; j < gh; j++ {
			ms[i] = append(ms[i], mapStrWalkable)
		}
	}
	for _, v := range g.obs {
		ms[int(v.i)][int(v.j)] = mapStrObs
	}
	ms[int(g.start.i)][int(g.start.j)] = mapStrStart
	ms[int(g.end.i)][int(g.end.j)] = mapStrEnd
	if g.state == mapStatusPause {
		g.OnStart()
	}
	g.state = oldState
	for _, v := range ms {
		_, err := f.Write(append(v, '\n'))
		if err != nil {
			panic(err)
		}
	}

}

func (g *Map) checkMapFile(f fyne.URIReadCloser) (res []string, err error) {
	reader := bufio.NewReader(f)
	defer f.Close()
	h := 0
	for {
		str, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		row := string(str)
		if h == 0 {
			h = len(row)
		} else if h != len(row) {
			return res, errors.New("invalid map height size,load fail")
		}
		res = append(res, row)
	}
	if len(res) < 1 {
		return res, errors.New("invalid map width size <1,load fail")
	}
	return res, nil
}

func (g *Map) Load(f fyne.URIReadCloser) {
	if g.state == mapStatusDumpOrLoading {
		return
	}
	g.OnClear()
	g.state = mapStatusDumpOrLoading
	mapStr, err := g.checkMapFile(f)
	if err != nil {
		dialog.ShowError(err, g.win)
		return
	}
	g.resize(len(mapStr), len(mapStr[0]), false)
	var (
		obs   []*grid
		start *grid
		end   *grid
	)
	for i, rows := range mapStr {
		for j, col := range rows {
			gg := g.grids[i][j]
			switch col {
			case mapStrObs:
				obs = append(obs, gg)
			case mapStrStart:
				start = gg
			case mapStrEnd:
				end = gg
			default:
				gg.SetWalkAble(true)
			}
		}
	}
	if start == nil || end == nil || start == end {
		g.initStartEnd()
		slog.Info("found start or end ==nil,init default")
	} else {
		g.SetStart(int(start.i), int(g.start.j))
		g.SetEnd(int(end.i), int(end.j))
	}
	g.obs = obs
	g.Refresh()
	g.state = mapStatusNone
}
