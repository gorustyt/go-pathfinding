package base

type LikeList struct {
	size int
	root *Node
}

func (list *LikeList) GetRoot() (root *Node) {
	return list.root
}

func (list *LikeList) Add(index int, Element interface{}) {
	if index == 0 {
		list.root.Element = Element
		list.size++
		return
	}
	newNode := NewNode(Element)
	node := list.findNodeByIndex(index - 1)
	if node == nil {
		list.root.Next = newNode
	} else {
		newNode.Next = node.Next
		node.Next = newNode
	}
	list.size++
}

func (list *LikeList) Get(index int) (Element interface{}) {
	node := list.findNodeByIndex(index)
	if node == nil {
		return node
	}
	return node.Element
}

func (list *LikeList) findNodeByIndex(index int) (node *Node) {
	if index >= list.size {
		return nil
	}
	head := list.root
	for i := 0; i < index; i++ {
		head = head.Next
	}
	return head
}

func (list *LikeList) Clear() {
	list.size = 0
	list.root = nil
}

func (list *LikeList) Remove(index int) {
	if index == 0 {
		list.Clear()
		return
	}
	node := list.findNodeByIndex(index - 1)
	if node == nil {
		return
	}
	node.Next = node.Next.Next
}

func (list *LikeList) Size() int {
	return list.size
}

func (list *LikeList) IsEmpty() bool {
	return list.Size() == 0
}

type Node struct {
	Element interface{}
	Next    *Node
}

func NewLinkList() *LikeList {
	return &LikeList{
		root: NewNode(nil),
	}
}
func NewNode(element interface{}) *Node {
	return &Node{
		Element: element,
	}
}
