package path_finding

import "math"

func (grid *Grid) JPFMoveDiagonallyIfAtMostOneObstacleFind(node *GridNodeInfo) (neighbors []*GridNode) {
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
			if grid.isWalkableAt(x, y+dy) || grid.isWalkableAt(x+dx, y) {
				neighbors = append(neighbors, grid.getNodeAt(x+dx, y+dy))

			}
			if !grid.isWalkableAt(x-dx, y) && grid.isWalkableAt(x, y+dy) {
				neighbors = append(neighbors, grid.getNodeAt(x-dx, y+dy))

			}
			if !grid.isWalkableAt(x, y-dy) && grid.isWalkableAt(x+dx, y) {
				neighbors = append(neighbors, grid.getNodeAt(x+dx, y-dy))
			}
		} else { // search horizontally/vertically
			if dx == 0 {
				if grid.isWalkableAt(x, y+dy) {
					neighbors = append(neighbors, grid.getNodeAt(x, y+dy))
					if !grid.isWalkableAt(x+1, y) {
						neighbors = append(neighbors, grid.getNodeAt(x+1, y+dy))
					}
					if !grid.isWalkableAt(x-1, y) {
						neighbors = append(neighbors, grid.getNodeAt(x-1, y+dy))
					}
				}
			} else {
				if grid.isWalkableAt(x+dx, y) {
					neighbors = append(neighbors, grid.getNodeAt(x+dx, y))

					if !grid.isWalkableAt(x, y+1) {
						neighbors = append(neighbors, grid.getNodeAt(x+dx, y+1))

					}
					if !grid.isWalkableAt(x, y-1) {
						neighbors = append(neighbors, grid.getNodeAt(x+dx, y-1))

					}
				}
			}
		}
	} else { // return all neighbors
		neighbors = grid.getNeighbors(node.GridNode, DiagonalMovementIfAtMostOneObstacle)
		for _, neighbor := range neighbors {
			neighbors = append(neighbors, grid.getNodeAt(neighbor.X, neighbor.Y))
		}
	}

	return neighbors
}

func (grid *Grid) JPFMoveDiagonallyIfAtMostOneObstacleJump(x, y, px, py int, endNode *GridNode) (jumpPoint *GridNode) {

	dx := x - px
	dy := y - py

	if !grid.isWalkableAt(x, y) {
		return nil
	}

	if grid.getNodeAt(x, y) == endNode {
		return grid.getNodeAt(x, y)
	}

	// check for forced neighbors
	// along the diagonal
	if dx != 0 && dy != 0 {
		if (grid.isWalkableAt(x-dx, y+dy) && !grid.isWalkableAt(x-dx, y)) ||
			(grid.isWalkableAt(x+dx, y-dy) && !grid.isWalkableAt(x, y-dy)) {
			return grid.getNodeAt(x, y)
		}
		// when moving diagonally, must check for vertical/horizontal jump points
		if grid.JPFMoveDiagonallyIfAtMostOneObstacleJump(x+dx, y, x, y, endNode) != nil || grid.JPFMoveDiagonallyIfAtMostOneObstacleJump(x, y+dy, x, y, endNode) != nil {
			return grid.getNodeAt(x, y)
		}
	}
	// horizontally/vertically else {
	if dx != 0 { // moving along x
		if (grid.isWalkableAt(x+dx, y+1) && !grid.isWalkableAt(x, y+1)) ||
			(grid.isWalkableAt(x+dx, y-1) && !grid.isWalkableAt(x, y-1)) {
			return grid.getNodeAt(x, y)
		}
	} else {
		if (grid.isWalkableAt(x+1, y+dy) && !grid.isWalkableAt(x+1, y)) ||
			(grid.isWalkableAt(x-1, y+dy) && !grid.isWalkableAt(x-1, y)) {
			return grid.getNodeAt(x, y)
		}
	}

	// moving diagonally, must make sure one of the vertical/horizontal
	// neighbors is open to allow the path
	if grid.isWalkableAt(x+dx, y) || grid.isWalkableAt(x, y+dy) {
		return grid.JPFMoveDiagonallyIfAtMostOneObstacleJump(x+dx, y+dy, x, y, endNode)
	} else {
		return nil
	}
	return jumpPoint
}
