package path_finding

import "github.com/lirongyangtao/mygo/base"

func (grid *Grid) PathFindingBreadthFirst(startX, startY, endX, endY int) (res []*GridNodeInfo) {
	openList := base.NewSimpleQueue()
	startNode := grid.getNodeAt(startX, startY)
	openList.Offer(startNode.ToGridNodeInfo())
	gridNodeInfo := map[*GridNode]*GridNodeInfo{}
	endNode := grid.getNodeAt(endX, endY)
	closed := map[*GridNodeInfo]struct{}{}
	for openList.Len() > 0 {
		node := openList.Poll().(*GridNodeInfo)
		if node.GridNode == endNode {
			return node.GetPaths()
		}

		neighbors := grid.getNeighbors(node.GridNode, grid.Config.DiagonalMovement)
		for _, neighbor := range neighbors {
			info, ok := gridNodeInfo[neighbor]
			if !ok {
				info = neighbor.ToGridNodeInfo()
				gridNodeInfo[neighbor] = info
			}
			if _, ok := closed[info]; ok {
				continue
			}
			info.Parent = node
			openList.Offer(info)
			closed[info] = struct{}{}
		}
	}
	return res
}
