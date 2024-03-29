package path_finding

import "math"

func (grid *Grid) JPFMoveDiagonallyIfNoObstaclesFind(node *GridNodeInfo) (neighbors []*GridNode) {
	var (
		parent         = node.Parent
		x              = node.X
		y              = node.Y
		px, py, dx, dy int
	)

	// directed pruning: can ignore most neighbors, unless forced.
	if parent != nil {
		px = parent.X
		py = parent.Y
		// get the normalized direction of travel
		dx = (x - px) / int(math.Max(math.Abs(float64(x-px)), float64(1)))
		dy = (y - py) / int(math.Max(math.Abs(float64(y-py)), float64(1)))

		// search diagonally
		if dx != 0 && dy != 0 {
			if grid.isWalkableAt(x, y+dy) {
				neighbors = append(neighbors, grid.getNodeAt(x, y+dy))
			}
			if grid.isWalkableAt(x+dx, y) {
				neighbors = append(neighbors, grid.getNodeAt(x+dx, y))
			}
			if grid.isWalkableAt(x, y+dy) && grid.isWalkableAt(x+dx, y) {
				neighbors = append(neighbors, grid.getNodeAt(x+dx, y+dy))
			}
		} else { // search horizontally/vertically
			var isNextWalkable bool
			if dx != 0 {
				isNextWalkable = grid.isWalkableAt(x+dx, y)
				var isTopWalkable = grid.isWalkableAt(x, y+1)
				var isBottomWalkable = grid.isWalkableAt(x, y-1)

				if isNextWalkable {
					neighbors = append(neighbors, grid.getNodeAt(x+dx, y))
					if isTopWalkable {
						neighbors = append(neighbors, grid.getNodeAt(x+dx, y+1))
					}
					if isBottomWalkable {
						neighbors = append(neighbors, grid.getNodeAt(x+dx, y-1))
					}
				}
				if isTopWalkable {
					neighbors = append(neighbors, grid.getNodeAt(x, y+1))
				}
				if isBottomWalkable {
					neighbors = append(neighbors, grid.getNodeAt(x, y-1))
				}
			} else if dy != 0 {
				isNextWalkable = grid.isWalkableAt(x, y+dy)
				var isRightWalkable = grid.isWalkableAt(x+1, y)
				var isLeftWalkable = grid.isWalkableAt(x-1, y)

				if isNextWalkable {
					neighbors = append(neighbors, grid.getNodeAt(x, y+dy))
					if isRightWalkable {
						neighbors = append(neighbors, grid.getNodeAt(x+1, y+dy))
					}
					if isLeftWalkable {
						neighbors = append(neighbors, grid.getNodeAt(x-1, y+dy))
					}
				}
				if isRightWalkable {
					neighbors = append(neighbors, grid.getNodeAt(x+1, y))

				}
				if isLeftWalkable {
					neighbors = append(neighbors, grid.getNodeAt(x-1, y))

				}
			}
		}
	}
	// return all neighbors else {
	neighborNodes := grid.getNeighbors(node.GridNode, DiagonalMovementOnlyWhenNoObstacles)
	for _, neighbor := range neighborNodes {
		neighbors = append(neighbors, grid.getNodeAt(neighbor.X, neighbor.Y))
	}

	return neighbors
}

func (grid *Grid) JPFMoveDiagonallyIfNoObstaclesJump(x, y, px, py int, endNode *GridNode) (jumpPoint *GridNode) {

	dx := x - px
	dy := y - py

	if !grid.isWalkableAt(x, y) {
		return nil
	}
	if !grid.TracePath(grid.getNodeAt(x, y)) {
		return
	}
	if grid.getNodeAt(x, y) == endNode {
		return grid.getNodeAt(x, y)
	}

	// check for forced neighbors
	// along the diagonal
	if dx != 0 && dy != 0 {
		// if ((grid.isWalkableAt(x - dx, y + dy) && !grid.isWalkableAt(x - dx, y)) ||
		// (grid.isWalkableAt(x + dx, y - dy) && !grid.isWalkableAt(x, y - dy))) {
		// return [x, y];
		// }
		// when moving diagonally, must check for vertical/horizontal jump points
		if grid.JPFMoveDiagonallyIfNoObstaclesJump(x+dx, y, x, y, endNode) != nil || grid.JPFMoveDiagonallyIfNoObstaclesJump(x, y+dy, x, y, endNode) != nil {
			return grid.getNodeAt(x, y)
		}
	} else { // horizontally/vertically
		if dx != 0 {
			if (grid.isWalkableAt(x, y-1) && !grid.isWalkableAt(x-dx, y-1)) ||
				(grid.isWalkableAt(x, y+1) && !grid.isWalkableAt(x-dx, y+1)) {
				return grid.getNodeAt(x, y)
			}
		} else if dy != 0 {
			if (grid.isWalkableAt(x-1, y) && !grid.isWalkableAt(x-1, y-dy)) ||
				(grid.isWalkableAt(x+1, y) && !grid.isWalkableAt(x+1, y-dy)) {
				return grid.getNodeAt(x, y)
			}
			// When moving vertically, must check for horizontal jump points
			// if (this._jump(x + 1, y, x, y) || this._jump(x - 1, y, x, y)) {
			// return [x, y];
			// }
		}
	}

	// moving diagonally, must make sure one of the vertical/horizontal
	// neighbors is open to allow the path
	if grid.isWalkableAt(x+dx, y) && grid.isWalkableAt(x, y+dy) {
		return grid.JPFMoveDiagonallyIfNoObstaclesJump(x+dx, y+dy, x, y, endNode)
	} else {
		return nil
	}
}
