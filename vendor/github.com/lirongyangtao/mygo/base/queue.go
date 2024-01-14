package base

import (
	"sync/atomic"
	"unsafe"
)

type Queue interface {
	Poll() (ele interface{})
	Offer(ele interface{})
	Len() int
	Empty() bool
}

type SimpleQueue struct {
	elements []interface{}
}

func NewSimpleQueue() Queue {
	return &SimpleQueue{}
}
func (q *SimpleQueue) Poll() (ele interface{}) {
	if q.Empty() {
		return nil
	}
	ele = q.elements[0]
	q.elements = q.elements[1:]
	return ele
}
func (q *SimpleQueue) Offer(ele interface{}) {
	q.elements = append(q.elements, ele)
}
func (q *SimpleQueue) Len() int {
	return len(q.elements)
}
func (q *SimpleQueue) Empty() bool {
	return q.Len() == 0
}

// 无锁队列
type lockFreeQueue struct {
	size       int64
	head, tail unsafe.Pointer
}

type queueNode struct {
	next  unsafe.Pointer
	Value interface{}
}

func NewQueue() Queue {
	stub := queueNode{}
	return &lockFreeQueue{
		head: unsafe.Pointer(&stub),
		tail: unsafe.Pointer(&stub),
	}
}

func (q *lockFreeQueue) Poll() (ele interface{}) {
	//head := q.Load(&q.head)
	return nil
}
func (q *lockFreeQueue) Offer(ele interface{}) {
	//newNode := queueNode{
	//	Value: ele,
	//}
	//preTail := q.Load(&q.tail)

}
func (q *lockFreeQueue) Len() int {
	return int(atomic.LoadInt64(&q.size))
}
func (q *lockFreeQueue) Empty() bool {
	return q.Len() == 0
}

func (q *lockFreeQueue) Load(p *unsafe.Pointer) *queueNode {
	return (*queueNode)(atomic.LoadPointer(p))
}
func (q *lockFreeQueue) Store(p *unsafe.Pointer, old, new *queueNode) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}
