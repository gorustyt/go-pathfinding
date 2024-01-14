package base

import (
	"fmt"
)

type EdgeInfo struct {
	From   interface{}
	To     interface{}
	Weight int
}

func (e *EdgeInfo) String() string {
	return fmt.Sprintf("%v==============>%v,weight=%v", e.From, e.To, e.Weight)
}

type GraphVisitor interface {
	Visitor(v *Vertex)
}
type Graph interface {
	AddEdge(from, to interface{})
	AddEdgeWithWeight(from, to interface{}, weight int)
	RemoveEdge(from, to interface{})
	RemoveVertex(v interface{})
	Bfs(v interface{}, visitor GraphVisitor)
	Dfs(v interface{}, visitor GraphVisitor)
	TopLogicSort() (res []interface{})
	Print()
	Prim(ele interface{}) []*EdgeInfo
	Kruskal() []*EdgeInfo
	Dijkstra(ele interface{}) map[interface{}]*MinPathInfo
	BellFord(ele interface{}) map[interface{}]*MinPathInfo
	Floyd() map[interface{}]map[interface{}]*MinPathInfo
}

type Vertex struct {
	Value    interface{}
	InEdges  map[string]*Edge
	OutEdges map[string]*Edge
}

func (v *Vertex) rangeOutEdges(fn func(edge *Edge)) {
	for _, v := range v.OutEdges {
		if fn != nil {
			fn(v)
		}
	}
}

func NewVertex(value interface{}) *Vertex {
	return &Vertex{
		Value:    value,
		InEdges:  map[string]*Edge{},
		OutEdges: map[string]*Edge{},
	}
}
func (v *Vertex) GetValue() interface{} {
	return v.Value
}

func (v *Vertex) getHashcode() string {
	return fmt.Sprintf("v:%v", v.Value)
}

func (v *Vertex) InEdgesContain(e *Edge) bool {
	if e == nil {
		return false
	}
	_, ok := v.InEdges[e.HashCode()]
	return ok
}

func (v *Vertex) OutEdgeContain(e *Edge) bool {
	if e == nil {
		return false
	}
	_, ok := v.OutEdges[e.HashCode()]
	return ok
}

func (v *Vertex) OutEdgesRemove(e *Edge) {
	delete(v.OutEdges, e.HashCode())
}
func (v *Vertex) InEdgesRemove(e *Edge) {
	delete(v.InEdges, e.HashCode())
}
func (v *Vertex) OutEdgesAdd(e *Edge) {
	v.OutEdges[e.HashCode()] = e
}
func (v *Vertex) InEdgesAdd(e *Edge) {
	v.InEdges[e.HashCode()] = e
}

type Edge struct {
	From   *Vertex
	To     *Vertex
	Weight int
}

func NewEdge(From, To *Vertex, weight int) (edge *Edge) {
	return &Edge{
		From:   From,
		To:     To,
		Weight: weight,
	}
}
func (edge *Edge) HashCode() string {
	return edge.From.getHashcode() + edge.To.getHashcode()
}

func (edge *Edge) Equal(other *Edge) bool {
	if other == nil {
		return false
	}
	return edge.HashCode() == other.HashCode()
}
func (edge *Edge) ToEdgeInfo() *EdgeInfo {
	return &EdgeInfo{
		From:   edge.From.Value,
		To:     edge.To.Value,
		Weight: edge.Weight,
	}
}

type graph struct {
	vertexes map[interface{}]*Vertex //顶点
	Edges    map[string]*Edge        //边
}

func NewGraph() Graph {
	return &graph{
		vertexes: map[interface{}]*Vertex{},
		Edges:    map[string]*Edge{},
	}
}
func (g *graph) GetVertexCount() int {
	return len(g.vertexes)
}

func (g *graph) AddVertex(value interface{}) *Vertex {
	v := NewVertex(value)
	g.vertexes[value] = v
	return v
}

func (g *graph) AddEdge(from, to interface{}) {
	g.AddEdgeWithWeight(from, to, 0)
}

func (g *graph) AddEdgeWithWeight(from, to interface{}, weight int) {
	fromV, ok := g.vertexes[from]
	if !ok {
		fromV = g.AddVertex(from)

	}
	toV, ok := g.vertexes[to]
	if !ok {
		toV = g.AddVertex(to)
	}
	edge := NewEdge(fromV, toV, weight)
	fromV.OutEdgesRemove(edge)
	toV.InEdgesRemove(edge)
	fromV.OutEdgesAdd(edge)
	toV.InEdgesAdd(edge)
	g.Edges[edge.HashCode()] = edge
}

func (g *graph) RemoveEdge(from, to interface{}) {
	fromV, ok := g.vertexes[from]
	if !ok {
		return
	}
	toV, ok := g.vertexes[to]
	if !ok {
		return
	}
	edge := NewEdge(fromV, toV, 0)
	fromV.OutEdgesRemove(edge)
	toV.InEdgesRemove(edge)
	delete(g.Edges, edge.HashCode())
}

func (g *graph) RemoveVertex(v interface{}) {
	V, ok := g.vertexes[v]
	if !ok {
		return
	}
	for _, v := range V.InEdges {
		v.From.OutEdgesRemove(v)
	}

	for _, v := range V.OutEdges {
		v.To.InEdgesRemove(v)
	}
	delete(g.vertexes, v)
}

// 广度搜索
func (g *graph) Bfs(v interface{}, visitor GraphVisitor) {
	V, ok := g.vertexes[v]
	if !ok {
		return
	}
	hasVisitor := map[*Vertex]struct{}{}
	queue := NewSimpleQueue()
	queue.Offer(V)
	hasVisitor[V] = struct{}{}
	for !queue.Empty() {
		ele := queue.Poll()
		V := ele.(*Vertex)
		if visitor != nil {
			visitor.Visitor(V)
		}
		for _, v := range V.OutEdges {
			if _, ok := hasVisitor[v.To]; ok {
				continue
			}
			queue.Offer(v.To)
			hasVisitor[v.To] = struct{}{}
		}
	}
}

// 深度搜索
func (g *graph) Dfs(v interface{}, visitor GraphVisitor) {
	V, ok := g.vertexes[v]
	if !ok {
		return
	}
	stack := NewSimpleStack()
	hasVisitor := map[*Vertex]struct{}{}
	stack.Push(V)
	if visitor != nil {
		visitor.Visitor(V)
	}
	hasVisitor[V] = struct{}{}
	for !stack.Empty() {
		ele := stack.Pop()
		V := ele.(*Vertex)
		for _, v := range V.OutEdges {
			if _, ok := hasVisitor[v.To]; ok {
				continue
			}
			stack.Push(v.From)
			stack.Push(v.To)
			if visitor != nil {
				visitor.Visitor(v.To)
			}
			hasVisitor[v.To] = struct{}{}
		}
	}
}

// 拓扑排序
func (g *graph) TopLogicSort() (res []interface{}) {
	queue := NewSimpleQueue()
	//记录每个顶点的入度是多少
	m := map[*Vertex]int{}
	for _, v := range g.vertexes {
		if len(v.InEdges) == 0 {
			queue.Offer(v)
			continue
		}
		m[v] = len(v.InEdges)
	}

	for !queue.Empty() {
		ele := queue.Poll()
		V := ele.(*Vertex)
		res = append(res, V.Value)
		for _, v := range V.OutEdges {
			m[v.To]--
			if m[v.To] == 0 {
				queue.Offer(v.To)
				delete(m, v.To)
				continue
			}
		}
	}
	return res
}
func (g *graph) toEdgeInfo(data *Edge) *EdgeInfo {
	return &EdgeInfo{
		From:   data.From.Value,
		To:     data.To.Value,
		Weight: data.Weight,
	}
}
func (g *graph) toEdgeInfoSlice(data []*Edge) (res []*EdgeInfo) {
	res = make([]*EdgeInfo, len(data))
	for index, v := range data {
		res[index] = g.toEdgeInfo(v)
	}
	return res
}

// prim 算法找最短路径
func (g *graph) Prim(ele interface{}) []*EdgeInfo {
	heap := NewQuadHeap(func(e1 interface{}, e2 interface{}) int32 {
		if e1.(*Edge).Weight < e2.(*Edge).Weight {
			return E1LessE2
		} else if e1.(*Edge).Weight > e2.(*Edge).Weight {
			return E1GenerateE2
		} else {
			return E1EqualE2
		}
	})
	var (
		edges   []*Edge
		visited = map[*Vertex]struct{}{}
	)
	v, ok := g.vertexes[ele]
	if !ok {
		return []*EdgeInfo{}
	}
	visited[v] = struct{}{}
	edgeSize := len(g.vertexes) - 1 //最小生成树边数量为顶点数量减一
	v.rangeOutEdges(func(edge *Edge) {
		heap.Add(edge)
	})
	for heap.Len() != 0 && len(edges) < edgeSize {
		edge := heap.Pop().(*Edge)
		if _, ok := visited[edge.To]; ok {
			continue
		}
		edges = append(edges, edge)
		visited[edge.To] = struct{}{}
		edge.To.rangeOutEdges(func(edge *Edge) {
			heap.Add(edge)
		})
	}
	return g.toEdgeInfoSlice(edges)
}
func (g *graph) Kruskal() []*EdgeInfo {
	var edges []*Edge
	heap := NewQuadHeap(func(e1 interface{}, e2 interface{}) int32 {
		if e1.(*Edge).Weight < e2.(*Edge).Weight {
			return E1LessE2
		} else if e1.(*Edge).Weight > e2.(*Edge).Weight {
			return E1GenerateE2
		} else {
			return E1EqualE2
		}
	})
	//将所有边加入heap
	for _, v := range g.vertexes {
		v.rangeOutEdges(func(edge *Edge) {
			heap.Add(edge)
		})
	}
	set := NewDisjointSet()
	edgeSize := len(g.vertexes) - 1 //最小生成树边数量为顶点数量减一
	for heap.Len() != 0 && len(edges) < edgeSize {
		edge := heap.Pop().(*Edge)
		//检查是否有环
		if set.IsSame(edge.From, edge.To) {
			continue
		}
		edges = append(edges, edge)
		set.Union(edge.From, edge.To)
	}

	return g.toEdgeInfoSlice(edges)
}

type MinPathInfo struct {
	Weight  int
	MinPath []*EdgeInfo
}

func (g *graph) finMinPath(pathTable map[*Vertex]*MinPathInfo) (res *Vertex, resWeight *MinPathInfo) {
	for k, weight := range pathTable {
		if res == nil {
			res = k
			resWeight = weight
			continue
		}
		if weight.Weight < resWeight.Weight {
			res = k
			resWeight = weight
		}
	}
	delete(pathTable, res)
	return
}

// 返回从给定顶点到所有路径的最短路径
func (g *graph) Dijkstra(ele interface{}) map[interface{}]*MinPathInfo {
	pathTable := map[*Vertex]*MinPathInfo{}           //key:顶点,value 权重，记
	selectPathTable := map[interface{}]*MinPathInfo{} //key:顶点,value 权重，记
	v, ok := g.vertexes[ele]
	if !ok {
		return selectPathTable
	}
	v.rangeOutEdges(func(edge *Edge) {
		path := &MinPathInfo{
			Weight:  edge.Weight,
			MinPath: []*EdgeInfo{edge.ToEdgeInfo()},
		}
		pathTable[edge.To] = path
	})

	for len(pathTable) > 0 {
		minVertex, minPath := g.finMinPath(pathTable)
		if _, ok := selectPathTable[minVertex.Value]; ok || minVertex == v {
			continue
		}
		selectPathTable[minVertex.Value] = minPath
		minVertex.rangeOutEdges(func(edge *Edge) {
			g.dijkstraRelax(pathTable, edge, minPath)
		})

	}
	return selectPathTable
}

func (g *graph) dijkstraRelax(pathTable map[*Vertex]*MinPathInfo,
	edge *Edge,
	minPath *MinPathInfo) {
	toWeight, ok := pathTable[edge.To]
	if !ok || (minPath.Weight+edge.Weight < toWeight.Weight) {
		path := &MinPathInfo{
			Weight: minPath.Weight + edge.Weight,
		}
		path.MinPath = append(minPath.MinPath, edge.ToEdgeInfo())
		pathTable[edge.To] = path
	}
}

func (g *graph) Floyd() map[interface{}]map[interface{}]*MinPathInfo {
	//初始化
	results := map[interface{}]map[interface{}]*MinPathInfo{} //key:顶点v1,value:map====>顶点v3，value 路径
	for _, v := range g.Edges {
		from, ok := results[v.From.Value]
		if !ok {
			from = map[interface{}]*MinPathInfo{}
		}
		path, ok := from[v.To.Value]
		if !ok {
			path = &MinPathInfo{}
		}
		path.Weight = v.Weight
		path.MinPath = append(path.MinPath, v.ToEdgeInfo())
		from[v.To.Value] = path
		results[v.From.Value] = from
	}
	getPathInfo := func(i, j interface{}) *MinPathInfo {
		from, ok := results[i]
		if !ok {
			return nil
		}
		path, ok := from[j]
		if !ok {
			return nil
		}
		return path
	}
	for k2 := range g.vertexes {
		for k1 := range g.vertexes {
			for k3 := range g.vertexes {

				//检查点重复
				if k1 == k2 || k1 == k3 || k2 == k3 {
					continue
				}
				//v1====>v2
				path1 := getPathInfo(k1, k2)
				if path1 == nil {
					continue
				}
				//v2====>v3
				path2 := getPathInfo(k2, k3)
				if path2 == nil {
					continue
				}
				//v1====>v3
				path3 := getPathInfo(k1, k3)
				//更新路径
				from, ok := results[k1]
				if !ok {
					from = map[interface{}]*MinPathInfo{}
				}
				if path3 != nil && path1.Weight+path2.Weight > path3.Weight {
					continue
					//比较大小
				}
				from[k3] = &MinPathInfo{
					Weight:  path1.Weight + path2.Weight,
					MinPath: append(path1.MinPath, path2.MinPath...),
				}
				results[k1] = from
			}
		}
	}
	return results
}

func (g *graph) BellFord(ele interface{}) map[interface{}]*MinPathInfo {
	selectPathTable := map[interface{}]*MinPathInfo{} //key:顶点,value 权重，记
	v, ok := g.vertexes[ele]
	if !ok {
		return nil
	}
	selectPathTable[v.Value] = &MinPathInfo{}
	for i := 0; i < len(g.vertexes)-1; i++ {
		for _, edge := range g.Edges {
			if minPath, ok := selectPathTable[edge.From.Value]; !ok {
				continue
			} else {
				g.bellFordRelax(selectPathTable, edge, minPath)
			}

		}
	}

	for i := 0; i < len(g.vertexes)-1; i++ {
		for _, v := range g.Edges {
			if minPath, ok := selectPathTable[v.From.Value]; !ok {
				continue
			} else {
				if g.bellFordRelax(selectPathTable, v, minPath) {
					fmt.Println("有负权边")
					return nil
				}
			}

		}
	}
	delete(selectPathTable, v.Value)
	return selectPathTable
}
func (g *graph) bellFordRelax(selectPathTable map[interface{}]*MinPathInfo,
	edge *Edge,
	minPath *MinPathInfo) bool {
	toWeight, ok := selectPathTable[edge.To.Value]
	if ok && (toWeight.Weight <= minPath.Weight+edge.Weight) {
		return false
	}
	path := &MinPathInfo{
		Weight: minPath.Weight + edge.Weight,
	}
	path.MinPath = append(minPath.MinPath, edge.ToEdgeInfo())
	selectPathTable[edge.To.Value] = path
	return true
}
func (g *graph) Print() {
	for _, v := range g.vertexes {
		for _, edge := range v.OutEdges {
			fmt.Println(fmt.Sprintf("[ from:%v=======>to:%v,  weight:%v]", edge.From.Value, edge.To.Value, edge.Weight))
		}
	}
}
