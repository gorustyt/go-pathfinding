package path_finding

import (
	"github.com/lirongyangtao/mygo/base"
	"math"
)

func (grid *Grid) PathFindingBiAStar(startX, startY, endX, endY int) (res []*GridNodeInfo) {
	cmp := func(e1 interface{}, e2 interface{}) int32 {
		if e1.(*GridNodeInfo).F > e2.(*GridNodeInfo).F {
			return base.E1GenerateE2
		} else if e1.(*GridNodeInfo).F < e2.(*GridNodeInfo).F {
			return base.E1LessE2
		} else {
			return base.E1EqualE2
		}
	}
	//开始=================================
	beginOpenList := base.NewQuadHeap(cmp)
	startNode := grid.getNodeAt(startX, startY)
	beginOpenList.Add(startNode.ToGridNodeInfo())
	beginGridNodeInfo := map[*GridNode]*GridNodeInfo{}
	beginClosed := map[*GridNode]struct{}{}
	//结束=================================
	EndOpenList := base.NewQuadHeap(cmp)

	endNode := grid.getNodeAt(endX, endY)
	EndOpenList.Add(endNode.ToGridNodeInfo())
	endGridNodeInfo := map[*GridNode]*GridNodeInfo{}
	EndClosed := map[*GridNode]struct{}{}

	for beginOpenList.Len() != 0 && EndOpenList.Len() != 0 {
		//从头部开始找
		nodeByStart := beginOpenList.Pop().(*GridNodeInfo)
		if nodeByStart.GridNode == endNode { //找到终点了
			return nodeByStart.GetPaths()
		}
		beginClosed[nodeByStart.GridNode] = struct{}{}
		neighbors := grid.getNeighbors(nodeByStart.GridNode, grid.Config.DiagonalMovement)
		for _, neighbor := range neighbors {
			if _, ok := beginClosed[neighbor]; ok { //该格子已经被访问过了
				continue
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
				info.H = grid.Config.weight * float64(grid.Config.Heuristic(
					int(math.Abs(float64(info.X)-float64(endX))),
					int(math.Abs(float64(info.Y)-float64(endY)))))
				info.F = info.G + info.H
				info.Parent = nodeByStart
			}
			if !info.Open {
				info.Open = true
				beginOpenList.Add(info)
			}

		}
		//从尾部开始找

		NodeByEnd := EndOpenList.Pop().(*GridNodeInfo)
		if NodeByEnd.GridNode == startNode { //找到终点了
			return NodeByEnd.GetPaths()
		}
		EndClosed[NodeByEnd.GridNode] = struct{}{}
		neighbors = grid.getNeighbors(NodeByEnd.GridNode, grid.Config.DiagonalMovement)
		for _, neighbor := range neighbors {
			if _, ok := EndClosed[neighbor]; ok { //该格子已经被访问过了
				continue
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
				info.H = grid.Config.weight * float64(grid.Config.Heuristic(
					int(math.Abs(float64(info.X)-float64(endX))),
					int(math.Abs(float64(info.Y)-float64(endY)))))
				info.F = info.G + info.H
				info.Parent = NodeByEnd
			}
			if !info.Open {
				info.Open = true
				EndOpenList.Add(info)
			}

		}

	}
	return res
}
