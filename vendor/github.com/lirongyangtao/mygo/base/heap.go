package base

type Heap interface {
	PeekTop() interface{}                      //查看堆顶元素
	Pop() interface{}                          //弹出堆顶元素
	Replace(ele interface{}) (pre interface{}) //替换堆顶元素
	Remove(index int)
	Len() int
	Add(ele ...interface{})
}

type binaryHeap struct {
	elements []interface{}
	cmp      CompareFunc //比较
}

func NewBinaryHeap(cmp CompareFunc) Heap {
	if cmp == nil {
		panic(any("cmp should not nil"))
	}
	return &binaryHeap{
		cmp: cmp,
	}
}

func NewQuadHeap(cmp CompareFunc) Heap {
	if cmp == nil {
		panic(any("cmp should not nil"))
	}
	return &QuadHeap{
		cmp: cmp,
	}
}
func (h *binaryHeap) Pop() interface{} { //弹出堆顶元素
	ele := h.PeekTop()
	if ele != nil {
		h.Remove(0)
	}
	return ele
}

func (h *binaryHeap) Add(eles ...interface{}) {
	for _, ele := range eles {
		h.elements = append(h.elements, ele)
		h.shiftUp(h.Len() - 1)
	}
}

func (h *binaryHeap) PeekTop() interface{} {
	if h.Len() <= 0 {
		return nil
	}
	return h.elements[0]
}

// 替换堆顶元素
func (h *binaryHeap) Replace(ele interface{}) (pre interface{}) {
	if h.Len() == 0 {
		h.elements[0] = ele
	} else {
		pre = h.elements[0]
		h.elements[0] = ele
		h.shiftDown(0)
	}
	return ele
}

func (h *binaryHeap) shiftDown(index int) {
	n := h.Len()
	if index < 0 || index >= n {
		return
	}
	ele := h.elements[index]
	for { //叶子节点数量
		left := index<<1 + 1
		if left >= n {
			break
		}
		w := h.elements[left]
		if left+1 < n && h.cmp(w, h.elements[left+1]) == E1GenerateE2 {
			w = h.elements[left+1]
			left++
		}
		if h.cmp(w, ele) == E1GenerateE2 {
			break
		}
		h.elements[index] = h.elements[left]
		index = left
	}
	h.elements[index] = ele

}

func binaryShiftDown(lessSwap *lessSwap, lo, hi, first int) {
	root := lo
	n := hi
	if root < 0 || root >= n {
		return
	}
	for { //叶子节点数量
		left := root<<1 + 1
		if left >= n {
			break
		}
		w := left
		if left+1 < n && lessSwap.Less(first+w, first+left+1) == E1GenerateE2 {
			w = left + 1
		}
		if lessSwap.Less(first+w, first+root) == E1GenerateE2 {
			break
		}
		lessSwap.Swap(first+root, first+w)
		root = first + w
	}
}

func (h *binaryHeap) shiftUp(index int) {
	if index < 0 || index >= h.Len() {
		return
	}
	ele := h.elements[index]
	for index > 0 {
		par := (index - 1) >> 1
		if h.cmp(ele, h.elements[par]) == E1GenerateE2 {
			break
		}
		h.elements[index] = h.elements[par]
		index = par
	}
	h.elements[index] = ele

}
func (h *binaryHeap) Remove(index int) {
	if index < 0 || index >= h.Len() {
		return
	}
	h.elements[index] = h.elements[h.Len()-1]
	h.elements[index] = -1
	h.shiftUp(index)

}
func (h *binaryHeap) Len() int {
	return len(h.elements)
}

type QuadHeap struct {
	cmp      CompareFunc //比较
	elements []any
}

func (h *QuadHeap) PeekTop() interface{} {
	if h.Len() <= 0 {
		return nil
	}
	return h.elements[0]
}
func (h *QuadHeap) Replace(ele interface{}) (pre interface{}) {
	if h.Len() == 0 {
		h.elements[0] = ele
	} else {
		pre = h.elements[0]
		h.elements[0] = ele
		h.shiftDown(0)
	}
	return ele
}

// 删除某个元素
func (h *QuadHeap) Remove(index int) {
	if index < 0 || index >= h.Len() {
		return
	}
	h.elements[index] = h.elements[h.Len()-1]
	h.elements = h.elements[:h.Len()-1]
	h.shiftDown(index)
}

func (h *QuadHeap) Len() int {
	return len(h.elements)
}

func (h *QuadHeap) Add(eles ...interface{}) {
	for _, ele := range eles {
		h.elements = append(h.elements, ele)
		h.shiftUp(len(h.elements) - 1)
	}
}
func (h *QuadHeap) shiftDown(index int) {
	n := h.Len()
	if index < 0 || index >= n {
		return
	}
	ele := h.elements[index]
	for { //叶子节点数量
		left := index<<2 + 1
		mid := left + 2
		if left >= n {
			break
		}
		w := h.elements[left]

		if left+1 < n && h.cmp(h.elements[left+1], w) == E1LessE2 {
			w = h.elements[left+1]
			left++
		}

		if mid < n {
			w3 := h.elements[mid]

			if mid+1 < n && h.cmp(h.elements[mid+1], w3) == E1LessE2 {
				w3 = h.elements[mid+1]
				mid++
			}

			if h.cmp(w3, w) == E1LessE2 {
				left = mid
				w = w3
			}
		}

		if h.cmp(w, ele) == E1GenerateE2 {
			break
		}
		h.elements[index] = h.elements[left]
		index = left

	}
	h.elements[index] = ele
}

// 下滤区间为[lo:hi]的值,通用接口
func quadShiftDown(lessSwap *lessSwap, lo, hi, first int) {
	root := lo
	if root < 0 || root >= hi {
		return
	}
	for { //叶子节点数量
		left := root<<2 + 1
		mid := left + 2
		if left >= hi {
			break
		}
		w := left

		if left+1 < hi && lessSwap.Less(first+left+1, first+left) == E1LessE2 {
			w = left + 1
			left++
		}

		if mid < hi {
			w3 := mid
			if mid+1 < hi && lessSwap.Less(first+mid+1, first+mid) == E1LessE2 {
				w3 = mid + 1
				mid++
			}

			if lessSwap.Less(first+w3, first+w) == E1LessE2 {
				left = mid
				w = w3
			}
		}

		if lessSwap.Less(first+w, first+root) == E1GenerateE2 {
			break
		}
		lessSwap.Swap(first+root, first+left)
		root = first + left
	}
}

func (h *QuadHeap) Pop() interface{} { //弹出堆顶元素
	ele := h.PeekTop()
	if ele != nil {
		h.Remove(0)
	}
	return ele
}
func (h *QuadHeap) shiftUp(index int) {
	if index < 0 || index >= h.Len() {
		return
	}
	ele := h.elements[index]

	for index > 0 {
		par := (index - 1) >> 2
		if h.cmp(ele, h.elements[par]) == E1GenerateE2 {
			break
		}
		h.elements[index] = h.elements[par]
		index = par
	}
	h.elements[index] = ele
}
