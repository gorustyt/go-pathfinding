package path_finding

func (grid *Grid) PathFindingBiDijkstra(startX, startY, endX, endY int) (res []*GridNodeInfo) {
	old := grid.Config.Heuristic
	grid.Config.Heuristic = func(x int, y int) int {
		return 0
	}
	res = grid.PathFindingBiAStar(startX, startY, endX, endY)
	grid.Config.Heuristic = old
	return res
}
