package path_finding

import (
	"math"
)

func (grid *Grid) PathFindingBiAStar(startX, startY, endX, endY int) (res []*PathPoint) {
	//开始=================================
	beginOpenList := newGridOpenList()
	startNode := grid.getNodeAt(startX, startY)
	beginOpenList.Push(startNode.ToGridNodeInfo())
	beginGridNodeInfo := map[*GridNode]*GridNodeInfo{}
	beginClosed := map[*GridNode]struct{}{}
	//结束=================================
	EndOpenList := newGridOpenList()

	endNode := grid.getNodeAt(endX, endY)
	EndOpenList.Push(endNode.ToGridNodeInfo())
	endGridNodeInfo := map[*GridNode]*GridNodeInfo{}
	EndClosed := map[*GridNode]struct{}{}

	for !beginOpenList.Empty() && !EndOpenList.Empty() {
		//从头部开始找
		nodeByStart := beginOpenList.Pop()
		if nodeByStart.GridNode == endNode { //找到终点了
			return nodeByStart.GetPaths()
		}
		beginClosed[nodeByStart.GridNode] = struct{}{}
		neighbors := grid.getNeighbors(nodeByStart.GridNode, grid.Config.DiagonalMovement)
		for _, neighbor := range neighbors {
			if _, ok := beginClosed[neighbor]; ok { //该格子已经被访问过了
				continue
			}
			if !grid.TracePath(neighbor) {
				return
			}
			ng := nodeByStart.G
			if neighbor.X-nodeByStart.X == 0 || neighbor.Y-nodeByStart.Y == 0 {
				ng = ng + 1
			} else {
				ng = ng + math.Sqrt2
			}
			info, ok := beginGridNodeInfo[neighbor]
			if !ok {
				info = neighbor.ToGridNodeInfo()
				beginGridNodeInfo[neighbor] = info
			}
			if !info.Open || ng < info.G {
				info.G = ng
				info.H = grid.Config.Weight * float64(grid.Config.Heuristic(
					int(math.Abs(float64(info.X)-float64(endX))),
					int(math.Abs(float64(info.Y)-float64(endY)))))
				info.F = info.G + info.H
				info.Parent = nodeByStart
			}
			if !info.Open {
				info.Open = true
				beginOpenList.Push(info)
			} else {
				beginOpenList.Update(info)
			}

		}
		//从尾部开始找

		NodeByEnd := EndOpenList.Pop()
		if NodeByEnd.GridNode == startNode { //找到终点了
			return NodeByEnd.GetPaths()
		}
		EndClosed[NodeByEnd.GridNode] = struct{}{}
		neighbors = grid.getNeighbors(NodeByEnd.GridNode, grid.Config.DiagonalMovement)
		for _, neighbor := range neighbors {
			if _, ok := EndClosed[neighbor]; ok { //该格子已经被访问过了
				continue
			}
			if !grid.TracePath(neighbor) {
				return
			}
			ng := NodeByEnd.G
			if neighbor.X-NodeByEnd.X == 0 || neighbor.Y-NodeByEnd.Y == 0 {
				ng = ng + 1
			} else {
				ng = ng + math.Sqrt2
			}
			info, ok := endGridNodeInfo[neighbor]
			if !ok {
				info = neighbor.ToGridNodeInfo()
				endGridNodeInfo[neighbor] = info
			}
			if !info.Open || ng < info.G {
				info.G = ng
				info.H = grid.Config.Weight * float64(grid.Config.Heuristic(
					int(math.Abs(float64(info.X)-float64(endX))),
					int(math.Abs(float64(info.Y)-float64(endY)))))
				info.F = info.G + info.H
				info.Parent = NodeByEnd
			}
			if !info.Open {
				info.Open = true
				EndOpenList.Push(info)
			} else {
				EndOpenList.Update(info)
			}

		}

	}
	return res
}
