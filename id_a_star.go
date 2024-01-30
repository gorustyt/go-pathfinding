package path_finding

import (
	"math"
	"time"
)

func (grid *Grid) idAStartSearch(heuristic func(a, b *GridNodeInfo) int, startTime int64, end *GridNodeInfo, node *GridNodeInfo, g float64, cutoff float64, nodeInfoGet func(node *GridNode) *GridNodeInfo) (*GridNodeInfo, float64) {
	if !grid.TracePath(node.GridNode) {
		return nil, 0
	}
	if grid.Config.IdAStarTimeLimit > 0 && time.Now().UnixMilli()-startTime > grid.Config.IdAStarTimeLimit*1000 {

		return nil, 0
	}
	f := g + float64(heuristic(node, end.ToGridNodeInfo()))*grid.Config.Weight
	if f > cutoff {
		return node, f
	}
	cost := func(a, b *GridNodeInfo) float64 {
		if a.X == b.X || a.Y == b.Y {
			return 1
		}
		return math.Sqrt2
	}
	if node.GridNode == end.GridNode {
		return node, f
	}

	minValue := math.MaxFloat64
	neighbours := grid.getNeighbors(node.GridNode, grid.Config.DiagonalMovement)
	for _, neighbour := range neighbours {
		findNode, t := grid.idAStartSearch(heuristic, startTime, end, nodeInfoGet(neighbour), g+cost(node, nodeInfoGet(neighbour)), cutoff, nodeInfoGet)
		if findNode == nil {
			continue
		}
		if findNode.GridNode == end.GridNode {
			nodeInfoGet(neighbour).Parent = node
			return findNode, t
		}
		if t < minValue {
			minValue = t
		}
	}
	if minValue == math.MaxFloat64 {
		panic(any("xxx not sure"))
	}
	return node, minValue
}
func (grid *Grid) PathFindingIdaStar(startX, startY, endX, endY int) (res []*PathPoint) {
	/**
	 * IDA* search implementation.
	 *
	 * @param {Node} The node currently expanding from.
	 * @param {number} Cost to reach the given node.
	 * @param {number} Maximum search depth (cut-off value).
	 * @param {Array<Array<number>>} The found route.
	 * @param {number} Recursion depth.
	 *
	 * @return {Object} either a number with the new optimal cut-off depth,
	 * or a valid node instance, in which case a path was found.
	 */
	// Node instance lookups:
	start := grid.getNodeAt(startX, startY)
	end := grid.getNodeAt(endX, endY)
	if start == nil || end == nil {
		return res
	}
	startTime := time.Now().UnixMilli()
	heuristic := func(a, b *GridNodeInfo) int {
		return grid.Config.Heuristic(int(math.Abs(float64(b.X)-float64(a.X))),
			int(math.Abs(float64(b.Y)-float64(a.Y))))
	}
	gridNodeInfo := map[*GridNode]*GridNodeInfo{}
	nodeInfoGet := func(node *GridNode) *GridNodeInfo {
		info, ok := gridNodeInfo[node]
		if !ok {
			info = node.ToGridNodeInfo()
			gridNodeInfo[node] = info
		}
		return info
	}
	cutOff := float64(heuristic(nodeInfoGet(start), nodeInfoGet(end)))
	var endNodeInfo *GridNodeInfo

	for {
		node, t := grid.idAStartSearch(heuristic, startTime, nodeInfoGet(end), nodeInfoGet(start), 0, cutOff, nodeInfoGet)
		if node == nil { //没有找到
			break
		}
		if node.GridNode == end {
			endNodeInfo = node
			break
		}
		cutOff = t
	}
	if endNodeInfo != nil {
		res = endNodeInfo.GetPaths()
		return res
	}
	return res
}
