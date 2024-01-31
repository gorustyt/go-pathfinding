package grid_map

import (
	"bufio"
	"fyne.io/fyne/v2"
)

func (g *Map) Dump(f fyne.URIWriteCloser) {

}

func (g *Map) Load(f fyne.URIReadCloser) {
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
	for i, rows := range mapStr {
		for j, col := range rows {
			grid := g.grids[i][j]
			switch col {
			case '1':
				g.obs = append(g.obs, g.grids[i][j])
			case '2':
				g.start = grid
			case '3':
				g.end = grid
			case '4':
				g.paths = append(g.paths, grid)
			}
		}
	}
	g.Refresh()
}
