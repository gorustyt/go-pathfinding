package path_finding

import (
	"github.com/lirongyangtao/mygo/base"
	"math"
)

func (grid *Grid) PathFindingJumpPoint(startX, startY, endX, endY int) (res []*GridNodeInfo) {
	type jumpFunc func(startX, startY, endX, endY int, endNode *GridNode) (jumpPoint *GridNode)
	type findNeighborsFunc func(node *GridNodeInfo) (neighbors []*GridNode)
	var jump jumpFunc = func(startX, startY, endX, endY int, endNode *GridNode) (jumpPoint *GridNode) { panic("no impl") }
	var findNeighbors findNeighborsFunc = func(node *GridNodeInfo) (neighbors []*GridNode) { panic("no impl") }
	if grid.Config.DiagonalMovement == DiagonalMovementNever {
		jump = grid.jPFNeverMoveDiagonallyJump
		findNeighbors = grid.jPFNeverMoveDiagonallyFind
	} else if grid.Config.DiagonalMovement == DiagonalMovementAlways {
		jump = grid.JPFAlwaysMoveDiagonallyJump
		findNeighbors = grid.JPFAlwaysMoveDiagonallyFind
	} else if grid.Config.DiagonalMovement == DiagonalMovementOnlyWhenNoObstacles {
		jump = grid.JPFMoveDiagonallyIfNoObstaclesJump
		findNeighbors = grid.JPFMoveDiagonallyIfNoObstaclesFind
	} else {
		jump = grid.JPFMoveDiagonallyIfAtMostOneObstacleJump
		findNeighbors = grid.JPFMoveDiagonallyIfAtMostOneObstacleFind
	}
	///=====================================================寻路代码=================================
	closed := map[*GridNode]struct{}{}
	gridNodeInfo := map[*GridNode]*GridNodeInfo{}
	heuristic := grid.Config.Heuristic
	startNode := grid.getNodeAt(startX, startY)
	endNode := grid.getNodeAt(endX, endY)
	openList := base.NewQuadHeap(func(e1 interface{}, e2 interface{}) int32 {
		if e1.(*GridNodeInfo).F > e2.(*GridNodeInfo).F {
			return base.E1GenerateE2
		} else if e1.(*GridNodeInfo).F < e2.(*GridNodeInfo).F {
			return base.E1LessE2
		} else {
			return base.E1EqualE2
		}
	})
	startInfo := startNode.ToGridNodeInfo()
	gridNodeInfo[startNode] = startInfo
	openList.Add(startInfo) //起点也是跳点

	identifySuccessors := func(node *GridNodeInfo) {
		x, y := node.X, node.Y
		neighbors := findNeighbors(node)
		for _, neighbor := range neighbors {
			jumpPoint := jump(neighbor.X, neighbor.Y, x, y, endNode)
			if jumpPoint == nil {
				continue
			}

			if _, ok := closed[jumpPoint]; ok {
				continue
			}
			jumpNode, ok := gridNodeInfo[jumpPoint]
			if !ok {
				jumpNode = jumpPoint.ToGridNodeInfo()
				gridNodeInfo[jumpPoint] = jumpNode
			}
			jx, jy := jumpNode.X, jumpNode.Y
			// include distance, as parent may not be immediately adjacent:
			d := octile(int(math.Abs(float64(jx-x))), int(math.Abs(float64(jy-y))))
			ng := node.G + float64(d) // next `g` value

			if !jumpNode.Open || ng < jumpNode.G {
				jumpNode.G = ng
				jumpNode.H = float64(heuristic(int(math.Abs(float64(jx-endX))), int(math.Abs(float64(jy-endY)))))
				jumpNode.F = jumpNode.G + jumpNode.H
				jumpNode.Parent = node
			}
			if !jumpNode.Open {
				jumpNode.Open = true
				openList.Add(jumpNode)
			}
		}

	}
	for openList.Len() != 0 {
		node := openList.Pop().(*GridNodeInfo)
		closed[node.GridNode] = struct{}{}

		if node.GridNode == endNode {
			return node.ExpandPath(node.GetPaths(), func(x, y int) *GridNodeInfo {
				pathNode := grid.getNodeAt(x, y)
				info := gridNodeInfo[pathNode]
				if info == nil {
					info = pathNode.ToGridNodeInfo()
				}
				return info
			})
		}

		identifySuccessors(node)
	}
	return res
}
