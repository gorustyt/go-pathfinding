package path_finding

import "testing"

func TestHeap(t *testing.T) {
	h := newGridOpenList()
	for i := 100; i >= 0; i-- {
		h.Push(&GridNodeInfo{F: float64(i)})
	}
	count := 0.
	for !h.Empty() {
		v := h.Pop()
		if count != v.F {
			t.Errorf("invalid real:%v expect:%v", v.F, count)
		}
		count++
	}
}
