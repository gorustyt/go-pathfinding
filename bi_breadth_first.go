package path_finding

func (grid *Grid) PathFindingBiBreadthFirst(startX, startY, endX, endY int) (res []*PathPoint) {
	//=====================================开始============================
	beginOpenList := newQueue()
	startNode := grid.getNodeAt(startX, startY)
	beginOpenList.PushBack(startNode.ToGridNodeInfo())
	beginGridNodeInfo := map[*GridNode]*GridNodeInfo{}

	beginClosed := map[*GridNodeInfo]struct{}{}
	//=====================================结束============================
	endOpenList := newQueue()
	endNode := grid.getNodeAt(endX, endY)
	endOpenList.PushBack(endNode.ToGridNodeInfo())

	endGridNodeInfo := map[*GridNode]*GridNodeInfo{}
	endClosed := map[*GridNodeInfo]struct{}{}

	for beginOpenList.Len() > 0 && endOpenList.Len() > 0 {
		nodeByBegin := beginOpenList.Front()
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
			beginOpenList.PushBack(info)
			beginClosed[info] = struct{}{}
		}

		nodeByEnd := endOpenList.Front()
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
			endOpenList.PushBack(info)
			endClosed[info] = struct{}{}
		}
	}
	return res
}
