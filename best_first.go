package path_finding

func (grid *Grid) PathFindingBestFirst(startX, startY, endX, endY int) (res []*PathPoint) {
	old := grid.Config.Heuristic
	grid.Config.Heuristic = func(x int, y int) int {
		return old(x, y) * 1000000
	}
	res = grid.PathFindingAStar(startX, startY, endX, endY)
	grid.Config.Heuristic = old
	return res
}
