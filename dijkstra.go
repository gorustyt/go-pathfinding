package path_finding

func (grid *Grid) PathFindingDijkstra(startX, startY, endX, endY int) (res []*GridNodeInfo) {
	old := grid.Config.Heuristic
	grid.Config.Heuristic = func(x int, y int) int {
		return 0
	}
	res = grid.PathFindingAStar(startX, startY, endX, endY)
	grid.Config.Heuristic = old
	return res
}
