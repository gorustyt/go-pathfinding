package base

type DisjointSet struct {
	nodes       map[interface{}]*DisjointNode
	nodesHeight map[*DisjointNode]int
}

type DisjointNode struct {
	Value  interface{}
	Parent *DisjointNode
}

func NewDisjointSet() *DisjointSet {
	return &DisjointSet{
		nodes:       map[interface{}]*DisjointNode{},
		nodesHeight: map[*DisjointNode]int{},
	}
}

func NewDisjointNode(value interface{}) *DisjointNode {
	node := &DisjointNode{
		Value: value,
	}
	//指向自己
	node.Parent = node
	return node
}
func (set *DisjointSet) find(value1 interface{}) *DisjointNode {
	node := set.nodes[value1]
	if node == nil {
		return nil
	}
	for node.Parent != node { //路径减半
		node.Parent = node.Parent.Parent
		node = node.Parent
	}
	return node
}

func (set *DisjointSet) Union(value1 interface{}, value2 interface{}) {
	p1 := set.find(value1)
	if p1 == nil {
		p1 = NewDisjointNode(value1)
		set.nodesHeight[p1] = 1
		set.nodes[value1] = p1
	}
	p2 := set.find(value2)
	if p2 == nil {
		p2 = NewDisjointNode(value2)
		set.nodesHeight[p2] = 1
		set.nodes[value2] = p2
	}
	if p1 == p2 {
		return
	}
	if set.nodesHeight[p1] > set.nodesHeight[p2] {
		p1.Parent = p2
	} else if set.nodesHeight[p1] < set.nodesHeight[p2] {
		p2.Parent = p1
	} else {
		p1.Parent = p2
		set.nodesHeight[p2]++
	}
}
func (set *DisjointSet) IsSame(value1 interface{}, value2 interface{}) bool {
	p1 := set.find(value1)
	p2 := set.find(value2)
	if p1 == nil && p2 == nil {
		return false
	}
	return p2 == p1
}

// 并查集
type DisjointSetInt struct {
	Set  []int
	Cap  int
	rank map[int]int //key 元素，value 元素所在高度
	size map[int]int //key 元素，value 元素所在元素个数
}

func NewDisjointIntSet(cap int) *DisjointSetInt {
	set := &DisjointSetInt{
		Set:  make([]int, cap),
		Cap:  cap,
		size: map[int]int{},
		rank: map[int]int{},
	}
	for index := range set.Set {
		set.Set[index] = index
		set.rank[index] = 1
		set.size[index] = 1
	}
	return set
}

// =======================================quickFind===================
func (set *DisjointSetInt) find(v int) int {
	p := set.Set[v]
	for p != set.Set[p] {
		p = set.Set[p]
	}
	return p
}

func (set *DisjointSetInt) checkRange(v ...int) {
	for _, value := range v {
		if value < 0 || value > set.Cap {
			panic(any("value not invalid"))
		}
	}
}
func (set *DisjointSetInt) checkRangeOk(v ...int) bool {
	for _, value := range v {
		if value < 0 || value > set.Cap {
			return false
		}
	}

	return true
}
func (set *DisjointSetInt) QuickFind_Union(v1 int, v2 int) {
	set.checkRange(v1, v2)
	p1 := set.find(v1)
	p2 := set.find(v2)
	if p1 == p2 {
		return
	}
	for index, _ := range set.Set {
		if index == p1 {
			set.Set[index] = p2
		}
	}
}
func (set *DisjointSetInt) IsSame(v1 int, v2 int) bool {
	if !set.checkRangeOk(v1) || !set.checkRangeOk(v2) {
		return false
	}
	return set.find(v1) == set.find(v2)
}

// =======================================quickUnion===================
// 没有优化
func (set *DisjointSetInt) QuickUnion_Union(v1 int, v2 int) {
	set.checkRange(v1, v2)
	p1 := set.find(v1)
	p2 := set.find(v2)
	if p1 == p2 {
		return
	}
	set.Set[p1] = p2
}

// 基于rank 优化
func (set *DisjointSetInt) QuickUnion_UnionPathRank(v1 int, v2 int) {
	set.checkRange(v1, v2)
	p1 := set.find(v1)
	p2 := set.find(v2)
	if p1 == p2 {
		return
	}
	if set.rank[p1] < set.rank[p2] {
		set.Set[p2] = p1
	} else if set.rank[p1] > set.rank[p2] {
		set.Set[p1] = p2
	} else {
		set.Set[p2] = p1
		set.rank[p1] += 1
	}

}

// 基于size优化
func (set *DisjointSetInt) QuickUnion_UnionPathSize(v1 int, v2 int) {
	set.checkRange(v1, v2)
	p1 := set.find(v1)
	p2 := set.find(v2)
	if p1 == p2 {
		return
	}
	if set.size[p1] < set.size[p2] {
		set.Set[p2] = p1
		set.size[p1] += set.size[p2]
	} else {
		set.Set[p1] = p2
		set.size[p2] += set.size[p1]
	}
}

// 路径压缩
func (set *DisjointSetInt) findCompress(v int) int {
	set.checkRange(v)
	if set.Set[v] != v {
		set.Set[v] = set.findCompress(set.Set[v])
	}
	return set.Set[v]
}

// 路径分裂
func (set *DisjointSetInt) findSplit(v int) int {
	set.checkRange(v)
	for set.Set[v] != v {
		p := set.Set[v]
		set.Set[v] = set.Set[p]
		v = p
	}
	return set.Set[v]
}

// 路径减半
func (set *DisjointSetInt) findHalving(v int) int {
	set.checkRange(v)
	for set.Set[v] != v {
		p := set.Set[v]
		set.Set[v] = set.Set[p]
		v = set.Set[v]
	}
	return set.Set[v]
}

func (set *DisjointSetInt) QuickUnion_UnionPathCompress(v1 int, v2 int) {
	set.checkRange(v1, v2)
	p1 := set.findCompress(v1)
	p2 := set.findCompress(v2)
	if p1 == p2 {
		return
	}
	if set.rank[p1] < set.rank[p2] {
		set.Set[p2] = p1
	} else if set.rank[p1] > set.rank[p2] {
		set.Set[p1] = p2
	} else {
		set.Set[p2] = p1
		set.rank[p1] += 1
	}
}

// 路径减半
func (set *DisjointSetInt) QuickUnion_UnionPathSplit(v1 int, v2 int) {
	set.checkRange(v1, v2)
	p1 := set.findSplit(v1)
	p2 := set.findSplit(v2)
	if p1 == p2 {
		return
	}
	if set.rank[p1] < set.rank[p2] {
		set.Set[p2] = p1
	} else if set.rank[p1] > set.rank[p2] {
		set.Set[p1] = p2
	} else {
		set.Set[p2] = p1
		set.rank[p1] += 1
	}
}

// 路径减半
func (set *DisjointSetInt) QuickUnion_UnionPathHalving(v1 int, v2 int) {
	set.checkRange(v1, v2)
	p1 := set.findHalving(v1)
	p2 := set.findHalving(v2)
	if p1 == p2 {
		return
	}
	if set.rank[p1] < set.rank[p2] {
		set.Set[p2] = p1
	} else if set.rank[p1] > set.rank[p2] {
		set.Set[p1] = p2
	} else {
		set.Set[p2] = p1
		set.rank[p1] += 1
	}
}
