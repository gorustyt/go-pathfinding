package path_finding

func (grid *Grid) PathFindingBiBestFirst(startX, startY, endX, endY int) (res []*GridNodeInfo) {
	old := grid.Config.Heuristic
	grid.Config.Heuristic = func(x int, y int) int {
		return old(x, y) * 1000000
	}
	res = grid.PathFindingBiAStar(startX, startY, endX, endY)
	grid.Config.Heuristic = old
	return res
}
