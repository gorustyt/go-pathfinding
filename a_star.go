package path_finding

import (
	"math"
)

// A*算法
func (grid *Grid) PathFindingAStar(startX, startY, endX, endY int) (res []*PathPoint) {
	openList := newGridOpenList()
	closed := map[*GridNode]struct{}{}
	gridNodeInfo := map[*GridNode]*GridNodeInfo{}
	startNode := grid.getNodeAt(startX, startY)
	openList.Push(startNode.ToGridNodeInfo())
	endNode := grid.getNodeAt(endX, endY)
	for !openList.Empty() {
		node := openList.Pop()
		if node.GridNode == endNode { //找到终点了
			return node.GetPaths()
		}
		closed[node.GridNode] = struct{}{}
		neighbors := grid.getNeighbors(node.GridNode, grid.Config.DiagonalMovement)
		for _, neighbor := range neighbors {
			if _, ok := closed[neighbor]; ok { //该格子已经被访问过了
				continue
			}
			if !grid.TracePath(neighbor) {
				return
			}
			ng := node.G
			if neighbor.X-node.X == 0 || neighbor.Y-node.Y == 0 {
				ng = ng + 1
			} else {
				ng = ng + math.Sqrt2
			}
			info, ok := gridNodeInfo[neighbor]
			if !ok {
				info = neighbor.ToGridNodeInfo()
				gridNodeInfo[neighbor] = info
			}
			if !info.Open || ng < info.G {
				info.G = ng
				info.H = grid.Config.Weight * float64(grid.Config.Heuristic(
					int(math.Abs(float64(info.X)-float64(endX))),
					int(math.Abs(float64(info.Y)-float64(endY)))))
				info.F = info.G + info.H
				info.Parent = node
			}
			if !info.Open {
				info.Open = true
				openList.Push(info)
			} else {
				openList.Update(info)
			}

		}

	}
	return nil
}
