package base

import "fmt"

type DoubleLinkList interface {
	Pop() *Element
	Front() *Element
	Back() *Element
	Len() int
	Remove(e *Element) any
	PushFront(v any) *Element
	PushBack(v any) *Element
	InsertBefore(v any, mark *Element) *Element
	InsertAfter(v any, mark *Element) *Element
	MoveToFront(e *Element)
	MoveToBack(e *Element)
	MoveBefore(e, mark *Element)
	MoveAfter(e, mark *Element)
	PushBackList(other DoubleLinkList)
	PushFrontList(other DoubleLinkList)
	Purge()
	Print()
}
type doubleLinkList struct {
	root Element
	len  int
}

func NewDoubleLinkList() DoubleLinkList {
	list := &doubleLinkList{}
	list.init()
	return list
}

type Element struct {
	next, prev *Element
	list       DoubleLinkList
	Value      any
}

func (e *Element) Next() *Element {
	return e.next
}
func (e *Element) Front() *Element {
	return e.next
}

func (list *doubleLinkList) Front() *Element {
	return list.root.next
}
func (list *doubleLinkList) Back() *Element {
	if list.Len() == 0 {
		return nil
	}
	return list.root.prev
}
func (list *doubleLinkList) init() {
	list.root.next = &list.root
	list.root.prev = &list.root
	list.len = 0
}

func (list *doubleLinkList) Purge() {
	list.init()
}
func (list *doubleLinkList) Len() int {
	return list.len
}
func (list *doubleLinkList) Remove(e *Element) any {
	if e == nil {
		return nil
	}
	list.remove(e)
	return e
}

func (list *doubleLinkList) PushFront(v any) *Element {
	return list.insertValue(v, &list.root)
}

func (list *doubleLinkList) PushBack(v any) *Element {
	return list.insertValue(v, list.root.prev)
}

func (list *doubleLinkList) InsertBefore(v any, mark *Element) *Element {
	if mark == nil {
		return nil
	}
	return list.insertValue(v, mark.prev)
}
func (list *doubleLinkList) InsertAfter(v any, mark *Element) *Element {
	if mark == nil {
		return nil
	}
	return list.insertValue(v, mark.next)
}
func (list *doubleLinkList) MoveToFront(e *Element) {
	if e.list != list || list.root.next == e {
		return
	}
	list.move(e, &list.root)
}
func (list *doubleLinkList) MoveToBack(e *Element) {
	if e.list != list || list.root.next == e {
		return
	}
	list.move(e, list.Back())
}
func (list *doubleLinkList) MoveBefore(e, mark *Element) {
	if e.list != list || mark.list != list || e == mark {
		return
	}
	if e == nil || mark == nil {
		return
	}
	list.move(e, mark.prev)
}
func (list *doubleLinkList) MoveAfter(e, mark *Element) {
	if e.list != list || mark.list != list || e == mark {
		return
	}
	if e == nil || mark == nil {
		return
	}
	list.move(e, mark.next)

}
func (list *doubleLinkList) PushBackList(other DoubleLinkList) {
	e := other.Front()
	for i := 0; i < other.Len(); i++ {
		list.PushBack(e.Value)
		e = e.next
	}
}
func (list *doubleLinkList) PushFrontList(other DoubleLinkList) {
	e := other.Front()
	for i := 0; i < other.Len(); i++ {
		list.PushFront(e.Value)
		e = e.next

	}
}

func (list *doubleLinkList) insertValue(v any, at *Element) *Element {
	return list.insert(&Element{Value: v}, at)
}

// 在at 位置插入e 这个节点
func (list *doubleLinkList) insert(e, at *Element) *Element {
	e.next = at.next
	e.prev = at

	at.next.prev = e
	at.next = e

	e.list = list
	list.len++
	return e
}
func (list *doubleLinkList) Pop() *Element {
	e := list.Front()
	if e == nil {
		return e
	}
	list.Remove(e)
	return e
}
func (list *doubleLinkList) remove(e *Element) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	list.len--
}

func (list *doubleLinkList) move(e, at *Element) {
	if e == at {
		return
	}
	//断掉e 之前的链接
	e.prev.next = e.next
	e.next.prev = e.prev

	e.next = at.next
	e.prev = at

	at.next.prev = e
	at.next = e
}

func (list *doubleLinkList) Print() {
	e := list.Front()
	for i := 0; i < list.Len(); i++ {
		fmt.Printf(" %v", e.Value)
		e = e.next

	}
}
