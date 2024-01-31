package grid_map

import (
	"bufio"
	"fyne.io/fyne/v2"
)

const (
	mapStrWalkable = '0'
	mapStrObs      = '1'
	mapStrStart    = '2'
	mapStrEnd      = '3'
)

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

func (g *Map) Load(f fyne.URIReadCloser) {
	if g.state == mapStatusDumpOrLoading {
		return
	}
	g.OnClear()
	g.state = mapStatusDumpOrLoading
	var mapStr []string
	reader := bufio.NewReader(f)
	defer f.Close()
	for {
		str, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		mapStr = append(mapStr, string(str))
	}
	startOk, EndOk := false, false
	for i, rows := range mapStr {
		for j, col := range rows {
			gg := g.grids[i][j]
			switch col {
			case mapStrObs:
				g.obs = append(g.obs, gg)
			case mapStrStart:
				g.start = gg
				startOk = true
			case mapStrEnd:
				EndOk = true
				g.end = gg
			default:
				gg.SetWalkAble(true)
			}
		}
	}
	if !startOk {
		g.start = g.grids[0][0]
	}
	if !EndOk {
		g.end = g.grids[1][1]
	}
	g.Refresh()
	g.state = mapStatusNone
}
