package path_finding

type queue struct {
	data []*GridNodeInfo
}

func newQueue() *queue {
	return &queue{}
}
func (q *queue) Len() int {
	return len(q.data)
}
func (q *queue) PushBack(v *GridNodeInfo) {
	q.data = append(q.data, v)
}

func (q *queue) Front() *GridNodeInfo {
	v := q.data[0]
	q.data = q.data[1:]
	return v
}

func (q *queue) Empty() bool {
	return q.Len() == 0
}
