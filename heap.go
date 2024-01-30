package path_finding

import "container/heap"

type gridOpenList struct {
	h *gridHeap
}

func newGridOpenList() *gridOpenList {
	return &gridOpenList{h: &gridHeap{}}
}

func (g *gridOpenList) Pop() *GridNodeInfo {
	return heap.Pop(g.h).(*GridNodeInfo)
}
func (g *gridOpenList) Update(v *GridNodeInfo) {
	heap.Fix(g.h, v.heapIndex)
}
func (g *gridOpenList) Push(v *GridNodeInfo) {
	heap.Push(g.h, v)
}

func (g *gridOpenList) Empty() bool {
	return g.h.Len() == 0
}

type gridHeap []*GridNodeInfo

func (g *gridHeap) Len() int {
	return len(*g)
}

func (g *gridHeap) Less(i, j int) bool {
	return (*g)[i].F < (*g)[j].F
}

func (g *gridHeap) Swap(i, j int) {
	(*g)[i], (*g)[j] = (*g)[j], (*g)[i]
	(*g)[i].heapIndex = i
	(*g)[j].heapIndex = j
}

func (g *gridHeap) Push(x any) {
	n := len(*g)
	item := x.(*GridNodeInfo)
	item.heapIndex = n
	*g = append(*g, item)
}

func (g *gridHeap) Pop() any {
	old := *g
	n := len(old)
	item := old[n-1]
	old[n-1] = nil      // avoid memory leak
	item.heapIndex = -1 // for safety
	*g = old[0 : n-1]
	return item
}
