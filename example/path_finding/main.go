package main

import (
	"fmt"
	pf "github.com/gorustyt/go-pathfinding"
)

func main() {
	w, h := 50, 25
	//创建地图
	m := pf.NewGrid(w, h)
	//设置不可行走
	m.SetWalkableAt(5, 5, false)
	m.SetWalkableAt(10, 10, false)
	//开始a*寻路
	m.PathFindingAStar(0, 0, w-1, h-1)
	//打印地图的寻路
	m.PathFindingPrint(pf.BiAStar, 0, 0, w-1, h-1)
	fmt.Println("====================================================")
	//打印地图的寻路
	m.PathFindingPrint(pf.AStar, 0, 0, w-1, h-1)
}
