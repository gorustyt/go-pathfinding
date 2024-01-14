package path_finding

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

const (
	DiagonalMovementNone                DiagonalMovement = 0
	DiagonalMovementAlways              DiagonalMovement = 1
	DiagonalMovementNever               DiagonalMovement = 2
	DiagonalMovementIfAtMostOneObstacle DiagonalMovement = 3
	DiagonalMovementOnlyWhenNoObstacles DiagonalMovement = 4
)

type DiagonalMovement int
type Heuristic func(x int, y int) int

// 一个格子存储的对象
type GridNodeInfo struct {
	*GridNode
	Parent *GridNodeInfo
	G      float64
	H      float64
	F      float64
	Open   bool //是否第一次访问
}

func (info *GridNodeInfo) GetPaths() (res []*GridNodeInfo) {
	node := info.Parent
	for node != nil {
		res = append(res, node)
		node = node.Parent
	}
	return
}
func (info *GridNodeInfo) ExpandPath(paths []*GridNodeInfo, m func(x, y int) *GridNodeInfo) (expanded []*GridNodeInfo) {

	var (
		length          = len(paths)
		coord0, coord1  *GridNodeInfo
		interpolated    [][]int
		interpolatedLen int
		i, j            int
	)
	if length < 2 {
		return expanded
	}
	for i = 0; i < length-1; i++ {
		coord0 = paths[i]
		coord1 = paths[i+1]

		interpolated = info.interpolate(coord0.X, coord0.Y, coord1.X, coord1.Y)
		interpolatedLen = len(interpolated)
		for j = 0; j < interpolatedLen-1; j++ {
			expanded = append(expanded, m(interpolated[j][0], interpolated[j][1]))
		}
	}
	expanded = append(expanded, paths[length-1])
	return expanded
}

func (info *GridNodeInfo) interpolate(x0, y0, x1, y1 int) (lines [][]int) {
	var (
		abs                     = func(x int) int { return int(math.Abs(float64(x))) }
		sx, sy, dx, dy, err, e2 int
	)

	dx = abs(x1 - x0)
	dy = abs(y1 - y0)
	sx = -1
	if x0 < x1 {
		sx = 1
	}
	sy = -1
	if y0 < y1 {
		sy = 1
	}
	err = dx - dy
	for {
		lines = append(lines, []int{x0, y0})
		if x0 == x1 && y0 == y1 {
			break
		}

		e2 = 2 * err
		if e2 > -dy {
			err = err - dy
			x0 = x0 + sx
		}
		if e2 < dx {
			err = err + dx
			y0 = y0 + sy
		}
	}
	return lines
}

// 一个格子存储的对象
type GridNode struct {
	X        int
	Y        int
	Walkable bool //是否可以行走
}

func (node *GridNode) ToGridNodeInfo() *GridNodeInfo {
	return &GridNodeInfo{
		GridNode: node,
	}
}

func NewNode(X int, Y int, Walkable bool) *GridNode {
	return &GridNode{
		X:        X,
		Y:        Y,
		Walkable: Walkable,
	}
}

// 地图
type Grid struct {
	Width  int
	Height int
	Nodes  map[int64]*GridNode //key :x|y,value:node
	Config *PathFindingConfig
}

func NewGrid(Width int, Height int, options ...PathFindingConfigOptions) *Grid {
	if Width <= 0 {
		Width = 100
	}
	if Height <= 0 {
		Height = 100
	}
	cfg := &PathFindingConfig{
		allowDiagonal: true,
	}
	for _, v := range options {
		v(cfg)
	}
	cfg.check()
	grid := &Grid{
		Width:  Width,
		Height: Height,
		Nodes:  map[int64]*GridNode{},
		Config: cfg,
	}
	grid.init()
	return grid
}

// 初始化地图
func (grid *Grid) init() {
	for i := 0; i < grid.Width; i++ {
		for j := 0; j < grid.Height; j++ {
			key := grid.packPos(i, j)
			grid.Nodes[key] = NewNode(i, j, true)
		}
	}

}

func (grid *Grid) SetMatrix(matrix [][]*GridNode) {

}

func (grid *Grid) packPos(x, y int) (pos int64) {
	return int64(x<<32 | y)
}

func (grid *Grid) unPackPos(pos int64) (x, y int) {
	x = int(pos >> 32)
	y = int(pos << 32 >> 32)
	return x, y
}
func (grid *Grid) SetWalkableAt(X int, Y int, Walkable bool) {
	key := grid.packPos(X, Y)
	node, ok := grid.Nodes[key]
	if !ok {
		return
	}
	node.Walkable = Walkable
	grid.Nodes[key] = node
}
func (grid *Grid) isInside(X int, Y int) bool {
	return (X >= 0 && X < grid.Width) && (Y >= 0 && Y < grid.Height)
}
func (grid *Grid) isWalkableAt(X int, Y int) bool {
	key := grid.packPos(X, Y)
	return grid.isInside(X, Y) && grid.Nodes[key].Walkable
}
func (grid *Grid) getNodeAt(X int, Y int) *GridNode {
	key := grid.packPos(X, Y)
	return grid.Nodes[key]
}

func (grid *Grid) Clone() *Grid {
	return nil
}
func (grid *Grid) PathFindingRoute(cmd PathFindingType) (cmdFunc PathFindingCmd) {
	wrap := func(f PathFindingCmd) PathFindingCmd {
		return func(startX, startY, endX, endY int) (res []*GridNodeInfo) {
			//begin := time.Now()
			res = f(startX, startY, endX, endY)
			//since := time.Since(begin)
			//fmt.Printf("%v: 耗时:%v\n", cmd, since)
			return res
		}
	}

	switch cmd {
	case AStar:
		cmdFunc = grid.PathFindingAStar
	case IdAStar:
		cmdFunc = grid.PathFindingIdaStar
	case Dijkstra:
		cmdFunc = grid.PathFindingDijkstra
	case BestFirst:
		cmdFunc = grid.PathFindingBestFirst
	case BreadthFirst:
		cmdFunc = grid.PathFindingBreadthFirst
	case JumpPoint:
		cmdFunc = grid.PathFindingJumpPoint
	case BiAStar:
		cmdFunc = grid.PathFindingBiAStar
	case BiBreadthFirst:
		cmdFunc = grid.PathFindingBiBreadthFirst
	case BiBestFirst:
		cmdFunc = grid.PathFindingBiBestFirst
	case BiDijkstra:
		cmdFunc = grid.PathFindingBiDijkstra
	default:
		panic("error cmd")
	}
	return wrap(cmdFunc)
}

/**
 * Get the neighbors of the given node.
 *
 *     offsets      diagonalOffsets:
 *  +---+---+---+    +---+---+---+
 *  |   | 0 |   |    | 0 |   | 1 |
 *  +---+---+---+    +---+---+---+
 *  | 3 |   | 1 |    |   |   |   |
 *  +---+---+---+    +---+---+---+
 *  |   | 2 |   |    | 3 |   | 2 |
 *  +---+---+---+    +---+---+---+
 *
 *  When allowDiagonal is true, if offsets[i] is valid, then
 *  diagonalOffsets[i] and
 *  diagonalOffsets[(i + 1) % 4] is valid.
 * @param {Node} node
 * @param {DiagonalMovement} diagonalMovement
 */
func (grid *Grid) getNeighbors(node *GridNode, diagonalMovement DiagonalMovement) (neighbors []*GridNode) {
	x := node.X
	y := node.Y

	s0 := false
	d0 := false
	s1 := false
	d1 := false
	s2 := false
	d2 := false
	s3 := false
	d3 := false
	// ↑
	if grid.isWalkableAt(x, y-1) {
		neighbors = append(neighbors, grid.Nodes[grid.packPos(x, y-1)])
		s0 = true
	}
	// →
	if grid.isWalkableAt(x+1, y) {
		neighbors = append(neighbors, grid.Nodes[grid.packPos(x+1, y)])
		s1 = true
	}
	// ↓
	if grid.isWalkableAt(x, y+1) {
		neighbors = append(neighbors, grid.Nodes[grid.packPos(x, y+1)])
		s2 = true
	}
	// ←
	if grid.isWalkableAt(x-1, y) {
		neighbors = append(neighbors, grid.Nodes[grid.packPos(x-1, y)])
		s3 = true
	}

	if diagonalMovement == DiagonalMovementNever {
		return neighbors
	}

	if diagonalMovement == DiagonalMovementOnlyWhenNoObstacles {
		d0 = s3 && s0
		d1 = s0 && s1
		d2 = s1 && s2
		d3 = s2 && s3
	} else if diagonalMovement == DiagonalMovementIfAtMostOneObstacle {
		d0 = s3 || s0
		d1 = s0 || s1
		d2 = s1 || s2
		d3 = s2 || s3
	} else if diagonalMovement == DiagonalMovementAlways {
		d0 = true
		d1 = true
		d2 = true
		d3 = true
	} else {
		panic(any("Incorrect value of diagonalMovement"))
	}

	// ↖
	if d0 && grid.isWalkableAt(x-1, y-1) {
		neighbors = append(neighbors, grid.Nodes[grid.packPos(x-1, y-1)])
	}
	// ↗
	if d1 && grid.isWalkableAt(x+1, y-1) {
		neighbors = append(neighbors, grid.Nodes[grid.packPos(x+1, y-1)])
	}
	// ↘
	if d2 && grid.isWalkableAt(x+1, y+1) {
		neighbors = append(neighbors, grid.Nodes[grid.packPos(x+1, y+1)])
	}
	// ↙
	if d3 && grid.isWalkableAt(x-1, y+1) {
		neighbors = append(neighbors, grid.Nodes[grid.packPos(x-1, y+1)])
	}

	return
}

func (grid *Grid) PathFindingPrint(cmd PathFindingType, startX, startY, endX, endY int) {
	arr := grid.getPrintMap()
	arr[startX][startY] = "\033[32ms\033[0m"
	arr[endX][endY] = "\033[31me\033[0m"
	paths := grid.PathFindingRoute(cmd)(startX, startY, endX, endY)
	for i := len(paths) - 2; i >= 0; i-- {
		path := paths[i]
		arr[path.X][path.Y] = "\u001B[34mw\u001B[0m"

	}
	grid.print(arr, true)
}

// 终端打印地图
func (grid *Grid) PrintMap() {
	arr := grid.getPrintMap()
	grid.print(arr)
}

func (grid *Grid) getPrintMap() (arr [][]string) {
	//初始化
	arr = make([][]string, grid.Width)
	for index := range arr {
		arr[index] = make([]string, grid.Height)
	}
	for k, v := range grid.Nodes {
		x, y := grid.unPackPos(k)
		if v.Walkable {
			arr[x][y] = "0"
		} else {
			arr[x][y] = "\033[1;33m1\033[0m"
		}
	}
	return
}

func (grid *Grid) printClear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
func (grid *Grid) print(arr [][]string, enablePrint ...bool) string {
	str := ""
	for i := 0; i < grid.Height; i++ {
		for j := 0; j < grid.Width; j++ {
			str += arr[j][i]
		}
		str += "\n"
	}
	if len(enablePrint) > 0 && enablePrint[0] {
		fmt.Printf("\r%v", str)
	}
	return str
}

func GridFromString(str string) (grid *Grid, startNode, endNode *GridNode) {
	rows := strings.Split(str, "\n")
	startX, startY := 0, 0
	endX, endY := 0, 0
	height := len(rows)
	width := len(rows[0])
	grid = NewGrid(width, height)
	for i, row := range rows {
		for j, v := range row {
			if v == '1' {
				grid.SetWalkableAt(j, i, false)
			} else if v == 'e' {
				endX, endY = j, i
			} else if v == 's' {
				startX, startY = j, i
			}
		}
	}
	return grid, grid.getNodeAt(startX, startY), grid.getNodeAt(endX, endY)
}
