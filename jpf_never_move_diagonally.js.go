package path_finding

import "math"

func (grid *Grid) jPFNeverMoveDiagonallyFind(node *GridNodeInfo) (neighbors []*GridNode) {
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

		if dx != 0 {
			if grid.isWalkableAt(x, y-1) {
				neighbors = append(neighbors, grid.getNodeAt(x, y-1))
			}
			if grid.isWalkableAt(x, y+1) {
				neighbors = append(neighbors, grid.getNodeAt(x, y+1))
			}
			if grid.isWalkableAt(x+dx, y) {
				neighbors = append(neighbors, grid.getNodeAt(x+dx, y))
			}
		} else if dy != 0 {
			if grid.isWalkableAt(x-1, y) {
				neighbors = append(neighbors, grid.getNodeAt(x-1, y))
			}
			if grid.isWalkableAt(x+1, y) {
				neighbors = append(neighbors, grid.getNodeAt(x+1, y))
			}
			if grid.isWalkableAt(x, y+dy) {
				neighbors = append(neighbors, grid.getNodeAt(x, y+dy))
			}
		}
	} else { // return all neighbors
		neighborNodes := grid.getNeighbors(node.GridNode, DiagonalMovementNever)
		for _, neighbor := range neighborNodes {
			neighbors = append(neighbors, grid.getNodeAt(neighbor.X, neighbor.Y))
		}
	}
	return neighbors
}

func (grid *Grid) jPFNeverMoveDiagonallyJump(x, y, px, py int, endNode *GridNode) (jumpPoint *GridNode) {
	dx := x - px
	dy := y - py

	if !grid.isWalkableAt(x, y) {
		return nil
	}
	if grid.getNodeAt(x, y) == endNode {
		return grid.getNodeAt(x, y)
	}

	if dx != 0 {
		if grid.isWalkableAt(x, y-1) && !grid.isWalkableAt(x-dx, y-1) ||
			grid.isWalkableAt(x, y+1) && !grid.isWalkableAt(x-dx, y+1) {
			return grid.getNodeAt(x, y)
		}
	} else if dy != 0 {
		if (grid.isWalkableAt(x-1, y) && !grid.isWalkableAt(x-1, y-dy)) ||
			(grid.isWalkableAt(x+1, y) && !grid.isWalkableAt(x+1, y-dy)) {
			return grid.getNodeAt(x, y)
		}
		//When moving vertically, must check for horizontal jump points
		if grid.jPFNeverMoveDiagonallyJump(x+1, y, x, y, endNode) != nil || grid.jPFNeverMoveDiagonallyJump(x-1, y, x, y, endNode) != nil {
			return grid.getNodeAt(x, y)
		}
	} else {
		panic("Only horizontal and vertical movements are allowed")
	}

	return grid.jPFNeverMoveDiagonallyJump(x+dx, y+dy, x, y, endNode)

}
