package path_finding

import "github.com/lirongyangtao/mygo/base"

func (grid *Grid) PathFindingBiBreadthFirst(startX, startY, endX, endY int) (res []*GridNodeInfo) {
	//=====================================开始============================
	beginOpenList := base.NewSimpleQueue()
	startNode := grid.getNodeAt(startX, startY)
	beginOpenList.Offer(startNode.ToGridNodeInfo())
	beginGridNodeInfo := map[*GridNode]*GridNodeInfo{}

	beginClosed := map[*GridNodeInfo]struct{}{}
	//=====================================结束============================
	endOpenList := base.NewSimpleQueue()
	endNode := grid.getNodeAt(endX, endY)
	endOpenList.Offer(endNode.ToGridNodeInfo())

	endGridNodeInfo := map[*GridNode]*GridNodeInfo{}
	endClosed := map[*GridNodeInfo]struct{}{}

	for beginOpenList.Len() > 0 && endOpenList.Len() > 0 {
		nodeByBegin := beginOpenList.Poll().(*GridNodeInfo)
		if nodeByBegin.GridNode == endNode {
			return nodeByBegin.GetPaths()
		}

		neighbors := grid.getNeighbors(nodeByBegin.GridNode, grid.Config.DiagonalMovement)
		for _, neighbor := range neighbors {
			info, ok := beginGridNodeInfo[neighbor]
			if !ok {
				info = neighbor.ToGridNodeInfo()
				beginGridNodeInfo[neighbor] = info
			}
			if _, ok := beginClosed[info]; ok {
				continue
			}
			info.Parent = nodeByBegin
			beginOpenList.Offer(info)
			beginClosed[info] = struct{}{}
		}

		nodeByEnd := endOpenList.Poll().(*GridNodeInfo)
		if nodeByEnd.GridNode == startNode {
			return nodeByEnd.GetPaths()
		}

		neighbors = grid.getNeighbors(nodeByEnd.GridNode, grid.Config.DiagonalMovement)
		for _, neighbor := range neighbors {
			info, ok := endGridNodeInfo[neighbor]
			if !ok {
				info = neighbor.ToGridNodeInfo()
				endGridNodeInfo[neighbor] = info
			}
			if _, ok := endClosed[info]; ok {
				continue
			}
			info.Parent = nodeByEnd
			endOpenList.Offer(info)
			endClosed[info] = struct{}{}
		}
	}
	return res
}
