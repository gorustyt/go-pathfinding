package base

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

type visitorWrap struct {
	NodeVisitor
	stop bool
}

func (visit *visitorWrap) Visit(node *BinaryTreeNode) bool {
	if visit.NodeVisitor.Visit(node) {
		visit.stop = true
	}
	return visit.stop
}
func (visit *visitorWrap) IsStop() bool {
	return visit.stop
}

type NodeVisitor interface {
	Visit(node *BinaryTreeNode) bool
}
type PrintNodeVisitor struct {
}

func (visit *PrintNodeVisitor) Visit(node *BinaryTreeNode) bool {
	fmt.Println(node.element)
	return false
}

// 二叉树
type BinaryTree struct {
	size int
	Root *BinaryTreeNode
	cmp  CompareFunc //比较
}

type BinaryTreeNode struct {
	parent   *BinaryTreeNode //父节点
	left     *BinaryTreeNode //左子节点
	right    *BinaryTreeNode //右子节点
	element  any             //元素
	tmpIndex int             //在二叉树中索引,打印辅助
	height   int             //高度在Avl树中使用
}

func (n *BinaryTreeNode) IsLeaf() bool {
	return n.left == nil && n.right == nil
}

// 是否有2个孩子
func (n *BinaryTreeNode) HasTwoChildren() bool {
	if n == nil {
		return false
	}
	return n.left != nil && n.right != nil
}
func (n *BinaryTreeNode) GetElement() interface{} {
	if n == nil {
		return nil
	}
	return n.element
}
func NewBinaryTreeNode(element any, parent *BinaryTreeNode) *BinaryTreeNode {
	return &BinaryTreeNode{
		parent:  parent,
		element: element,
		height:  1,
	}
}

func NewBinaryTree(cmp CompareFunc) *BinaryTree {
	return &BinaryTree{
		cmp: cmp,
	}
}

func (tree *BinaryTree) Len() int {
	return tree.size
}

// 树的高度
func (tree *BinaryTree) Height() int {
	return tree.heightByIter(tree.Root)
}

// 递归实现
func (tree *BinaryTree) height(node *BinaryTreeNode) int {
	if node == nil {
		return 0
	}
	return 1 + int(math.Max(float64(tree.height(node.left)), float64(tree.height(node.right))))
}

// 迭代实现
func (tree *BinaryTree) heightByIter(node *BinaryTreeNode) (height int) {
	if node == nil {
		return 0
	}
	que := NewSimpleQueue()
	que.Offer(node)
	size := que.Len()
	for que.Len() != 0 {
		node := que.Poll().(*BinaryTreeNode)
		size--
		if node.left != nil {
			que.Offer(node.left)
		}
		if node.right != nil {
			que.Offer(node.right)
		}
		if size == 0 {
			height++
			size = que.Len()
		}
	}
	return
}

// 获取二叉树的每一层节点索引
func (tree *BinaryTree) getTreeNodeIndexMap() (nodeIndexMap map[int]map[int]interface{}, height int) {
	node := tree.Root
	if node == nil {
		return nil, height
	}
	node.tmpIndex = 0
	que := NewSimpleQueue()
	que.Offer(node)
	size := que.Len()
	height = 0
	//获取元素最大字符长度
	maxLen := 0
	nodeIndexMap = map[int]map[int]interface{}{}
	for que.Len() != 0 {
		node = que.Poll().(*BinaryTreeNode)
		tmp, ok := nodeIndexMap[height]
		if !ok {
			tmp = map[int]interface{}{}
		}
		tmp[node.tmpIndex] = node.element
		nodeIndexMap[height] = tmp
		if len(fmt.Sprintf("%v", node.element)) > maxLen {
			maxLen = len(fmt.Sprintf("%v", node.element))
		}
		size--
		if node.left != nil {
			node.left.tmpIndex = 2*node.tmpIndex + 1
			que.Offer(node.left)
		}
		if node.right != nil {
			node.right.tmpIndex = 2*node.tmpIndex + 2
			que.Offer(node.right)
		}
		if size == 0 {
			height++
			size = que.Len()
		}
	}
	return nodeIndexMap, height
}

// //树状字符
// //"└" "─" "┌─────────┴─────────┐" "┬"
// //"├"
// //"─┴─"
func (tree *BinaryTree) TreePrint() {
	nodeIndexMap, height := tree.getTreeNodeIndexMap()
	if nodeIndexMap == nil {
		return
	}
	result := map[int][]string{} //height===>index:value
	//从最后一层开始生成索引
	prefix := strings.Repeat(" ", 10)
	for i := height - 1; i >= 0; i-- {
		indexMap := nodeIndexMap[i]
		realIndexSlice, ok := result[i]
		if !ok {
			realIndexSlice = make([]string, 1<<(i))
		}
		for index, value := range indexMap {
			realIndex := index - (1<<(i) - 1) //每一层从零开始排
			realIndexSlice[realIndex] = fmt.Sprintf("%v", value)
		}
		result[i] = realIndexSlice
	}
	type IndexNode struct {
		Left       int
		Mid        int
		Right      int
		LeftEmpty  bool
		RightEmpty bool
	}
	//生成字符串及索引
	strMap := map[int][]IndexNode{} //height========>
	strResult := map[int]string{}
	fixSpace := strings.Repeat(" ", 3)
	fixNode := strings.Repeat(" ", 3)
	getStr := func(height int) (string, string) {
		//生成最后一层
		indexMap := result[height]
		str := ""
		lineStr := ""
		node := IndexNode{}
		strTmp, ok := strMap[height+1]
		if !ok {
			str += prefix
			isBegin := false
			for _, value := range indexMap {
				valueStr := fmt.Sprintf("%v", value)
				valueStrLen := utf8.RuneCountInString(valueStr)
				if valueStr == "" {
					valueStrLen = utf8.RuneCountInString(fixNode)
				}

				if !isBegin {
					node.Left = utf8.RuneCountInString(str) + valueStrLen/2
					node.LeftEmpty = valueStr == ""
					isBegin = true
				} else {
					node.RightEmpty = valueStr == ""
					node.Right = utf8.RuneCountInString(str) + valueStrLen/2
					node.Mid = node.Left + (node.Right-node.Left)/2
					strMap[height] = append(strMap[height], node)
					node = IndexNode{}
					isBegin = false
				}
				if valueStr == "" {
					valueStr = fixNode
				}
				str += valueStr
				str += fixSpace
			}
		} else {
			isBegin := false
			begin := 0
			lineBegin := false
			for index, nodeTmp := range strTmp {
				value := indexMap[index]
				for i := begin; i <= nodeTmp.Right; i++ {
					if i == nodeTmp.Left {
						str += " "
						if nodeTmp.LeftEmpty {
							lineStr += " "
						} else {
							lineStr += "┌"
							lineBegin = true
						}
					} else if i == nodeTmp.Mid {
						//处理因为字符长度导致的偏移
						if len(str) > nodeTmp.Mid {
							str = str[:nodeTmp.Mid]
						}
						valueStr := fmt.Sprintf("%v", value)
						valueStrLen := utf8.RuneCountInString(valueStr)
						if valueStr == "" {
							valueStrLen = utf8.RuneCountInString(fixNode)
						}
						if !isBegin {
							node.LeftEmpty = valueStr == ""
							node.Left = utf8.RuneCountInString(str) + valueStrLen/2
							isBegin = true
						} else {
							node.RightEmpty = valueStr == ""
							node.Right = utf8.RuneCountInString(str) + valueStrLen/2
							node.Mid = node.Left + (node.Right-node.Left)/2
							strMap[height] = append(strMap[height], node)
							node = IndexNode{}
							isBegin = false
						}
						if valueStr == "" {
							valueStr = fixNode
						}
						str += valueStr
						if nodeTmp.LeftEmpty && nodeTmp.RightEmpty {
							lineStr += " "
						} else if nodeTmp.LeftEmpty {
							lineStr += "┴"
							lineBegin = true
						} else if nodeTmp.RightEmpty {
							lineStr += "┴"
							lineBegin = false
						} else {
							lineStr += "┴"
						}

					} else if i == nodeTmp.Right {
						str += " "
						if nodeTmp.RightEmpty {
							lineStr += " "
						} else {
							lineStr += "┐"

						}
						lineBegin = false
					} else {
						str += " "
						if lineBegin {
							lineStr += "─"
						} else {
							lineStr += " "
						}

					}

				}
				begin = nodeTmp.Right + 1
			}
		}

		return str, lineStr
	}

	//循环生成前面的索引

	for i := height - 1; i >= 0; i-- {
		strResult[i], strResult[i+1000] = getStr(i)
	}
	//打印
	for i := 0; i < height; i++ {
		fmt.Println(strResult[i])
		fmt.Println(strResult[i+1000])
	}
}

func (tree *BinaryTree) isLeaf(node *BinaryTreeNode) bool {
	return node.left == nil && node.right == nil
}
func (tree *BinaryTree) IsCompete() bool {
	node := tree.Root
	que := NewSimpleQueue()
	que.Offer(node)
	mustLeaf := false
	for que.Len() != 0 {
		node := que.Poll().(*BinaryTreeNode)
		if mustLeaf && !tree.isLeaf(node) {
			return false
		}
		if node.left != nil {
			que.Offer(node.left)
		} else if node.right != nil {
			return false
		}
		if node.right != nil {
			que.Offer(node.right)
		} else {
			mustLeaf = true
		}
	}
	return true
}

func (tree *BinaryTree) Add(ele any) {
	if tree.Root == nil {
		tree.Root = NewBinaryTreeNode(ele, nil)
		tree.size++
		return
	}
	node := tree.Root
	var parent *BinaryTreeNode
	var res int32
	for node != nil {
		parent = node
		res = tree.cmp(node.element, ele)
		if res == E1LessE2 {
			node = node.right
		} else if res == E1GenerateE2 {
			node = node.left
		} else {
			break
		}
	}
	if res == E1LessE2 {
		newNode := NewBinaryTreeNode(ele, parent)
		parent.right = newNode
		tree.size++
		tree.afterAdd(newNode)
	} else if res == E1GenerateE2 {
		newNode := NewBinaryTreeNode(ele, parent)
		parent.left = newNode
		tree.size++
		tree.afterAdd(newNode)
	} else {
		parent.element = ele
	}

}

func (tree *BinaryTree) PreOrder() {
	tree.preOrder(tree.Root, &visitorWrap{NodeVisitor: &PrintNodeVisitor{}})
}
func (tree *BinaryTree) visit(node *BinaryTreeNode, visitor *visitorWrap) (isStop bool) {
	if visitor != nil {
		return visitor.Visit(node)
	}
	return false
}

func (tree *BinaryTree) preOrder(node *BinaryTreeNode, visitor *visitorWrap) {
	if node == nil || visitor.IsStop() {
		return
	}
	if tree.visit(node, visitor) {
		return
	}
	tree.preOrder(node.left, visitor)
	tree.preOrder(node.right, visitor)
}
func (tree *BinaryTree) InOrder() {
	tree.inOrder(tree.Root, &visitorWrap{NodeVisitor: &PrintNodeVisitor{}})

}
func (tree *BinaryTree) inOrder(node *BinaryTreeNode, visitor *visitorWrap) {
	if node == nil || visitor.IsStop() {
		return
	}
	tree.inOrder(node.left, visitor)
	if tree.visit(node, visitor) {
		return
	}
	tree.inOrder(node.right, visitor)

}
func (tree *BinaryTree) PostOrder() {
	tree.postOrder(tree.Root, &visitorWrap{NodeVisitor: &PrintNodeVisitor{}})
}
func (tree *BinaryTree) postOrder(node *BinaryTreeNode, visitor *visitorWrap) {
	if node == nil || visitor.IsStop() {
		return
	}

	tree.postOrder(node.left, visitor)
	tree.postOrder(node.right, visitor)

}

func (tree *BinaryTree) LevelOrder() {
	tree.levelOrder(tree.Root, &visitorWrap{NodeVisitor: &PrintNodeVisitor{}})
}

func (tree *BinaryTree) levelOrder(node *BinaryTreeNode, visitor *visitorWrap) {
	if node == nil || visitor.IsStop() {
		return
	}
	que := NewSimpleQueue()
	que.Offer(node)
	for que.Len() != 0 {
		node := que.Poll().(*BinaryTreeNode)
		if tree.visit(node, visitor) {
			return
		}
		if node.left != nil {
			que.Offer(node.left)
		}
		if node.right != nil {
			que.Offer(node.right)
		}
	}
}

func (tree *BinaryTree) PreOrderByVisitor(visitor NodeVisitor) {
	tree.preOrder(tree.Root, &visitorWrap{NodeVisitor: visitor})
}
func (tree *BinaryTree) InOrderByVisitor(visitor NodeVisitor) {
	tree.inOrder(tree.Root, &visitorWrap{NodeVisitor: visitor})
}
func (tree *BinaryTree) PostOrderByVisitor(visitor NodeVisitor) {
	tree.postOrder(tree.Root, &visitorWrap{NodeVisitor: visitor})
}

func (tree *BinaryTree) LevelOrderByVisitor(visitor NodeVisitor) {
	tree.levelOrder(tree.Root, &visitorWrap{NodeVisitor: visitor})
}

// 找后继者
func (tree *BinaryTree) FindSuccessor(node *BinaryTreeNode) *BinaryTreeNode {
	if node == nil {
		return nil
	}
	if node.right != nil {
		node = node.right
		for node.left != nil {
			node = node.left
		}
		return node
	}

	for node.parent != nil {
		node = node.parent
		if node == node.parent.left {
			return node
		}
	}
	return node
}

// 根据元素找到节点
func (tree *BinaryTree) FindNode(ele interface{}) (*BinaryTreeNode, bool) {
	node := tree.Root
	for node != nil {
		if tree.cmp(ele, node.element) == E1GenerateE2 {
			node = node.right
		} else if tree.cmp(ele, node.element) == E1LessE2 {
			node = node.left
		} else {
			return node, true
		}
	}
	return nil, false
}

// 找前驱
func (tree *BinaryTree) Predecessor(node *BinaryTreeNode) *BinaryTreeNode {
	if node == nil {
		return nil
	}
	if node.left != nil {
		node = node.left
		for node.right != nil {
			node = node.right
		}
		return node
	}

	for node.parent != nil {
		node = node.parent
		if node == node.parent.right {
			return node
		}
	}
	return node
}

func (tree *BinaryTree) Remove(ele interface{}) {
	node, ok := tree.FindNode(ele)
	if !ok {
		return
	}
	if node.IsLeaf() { //是叶子节点
		tree.removeLeaf(node)
	} else if node.HasTwoChildren() { //度为2的节点
		replaceNode := tree.FindSuccessor(node)
		node.element = replaceNode.element
		if replaceNode.IsLeaf() { //是叶子节点
			tree.removeLeaf(replaceNode)
		} else {
			tree.remove1(replaceNode)
		}
	} else { //度为1的节点
		tree.remove1(node)
	}
	tree.size--
}

func (tree *BinaryTree) removeLeaf(node *BinaryTreeNode) {
	if node.parent == nil {
		tree.Root = nil
	} else if node == node.parent.left {
		node.parent.left = nil
	} else {
		node.parent.right = nil
	}
	tree.afterRemove(node)
}

// 移除度为1的节点
func (tree *BinaryTree) remove1(node *BinaryTreeNode) {
	if node.parent == nil { //移除根节点
		if node.left != nil {
			tree.Root = node.left
			node.left.parent = nil
		} else {
			tree.Root = node.right
			node.right.parent = nil
		}
		return
	}
	if node == node.parent.left {
		if node.left != nil {
			node.parent.left = node.left
			node.left.parent = node.parent
		} else {
			node.parent.left = node.right
			node.right.parent = node.parent
		}
	} else {
		if node.right != nil {
			node.parent.right = node.right
			node.right.parent = node.parent
		} else {
			node.parent.right = node.left
			node.left.parent = node.parent
		}
	}
	tree.afterRemove(node)
}

// Avl树需要的接口
func (tree *BinaryTree) updateHeight(node *BinaryTreeNode) {
	if node == nil {
		return
	}
	leftHeight := 0
	rightHeight := 0
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}
	node.height = int(math.Max(float64(leftHeight), float64(rightHeight))) + 1
}

// 恢复平衡,分开旋转
func (tree *BinaryTree) reBalance1(grand *BinaryTreeNode) {
	if grand == nil {
		return
	}
	parent := tree.getTallerChild(grand)
	node := tree.getTallerChild(parent)
	if node == nil {
		return
	}
	if parent == grand.left && node == parent.left { //LL
		tree.rotateLL(grand, parent)
	} else if parent == grand.left && node == parent.right { //LR
		tree.rotateRR(parent, node)
		tree.rotateLL(grand, node)
	}
	if parent == grand.right && node == parent.right { //RR
		tree.rotateRR(grand, parent)
	} else if parent == grand.right && node == parent.left { //RL
		tree.rotateLL(parent, node)
		tree.rotateRR(grand, node)
	}
	tree.updateHeight(grand)
	tree.updateHeight(parent)
	tree.updateHeight(node)
}
func (tree *BinaryTree) rotateLL(grand, parent *BinaryTreeNode) {
	if grand.parent != nil {
		if grand == grand.parent.left {
			grand.parent.left = parent
		}
		if grand == grand.parent.right {
			grand.parent.right = parent
		}
		parent.parent = grand.parent
	} else {
		tree.Root = parent
		parent.parent = nil
	}
	grand.left = parent.right
	if parent.right != nil {
		parent.right.parent = grand
	}

	parent.right = grand
	grand.parent = parent

}

func (tree *BinaryTree) rotateRR(grand, parent *BinaryTreeNode) {
	if grand.parent != nil {
		if grand == grand.parent.left {
			grand.parent.left = parent
		}
		if grand == grand.parent.right {
			grand.parent.right = parent
		}
		parent.parent = grand.parent
	} else {
		tree.Root = parent
		parent.parent = nil
	}
	grand.right = parent.left
	if parent.left != nil {
		parent.left.parent = grand
	}

	parent.left = grand
	grand.parent = parent

}

// 恢复平衡,统一旋转
func (tree *BinaryTree) reBalance2(node *BinaryTreeNode) {

}

// 统一旋转
func (tree *BinaryTree) rotate() {

}

func (tree *BinaryTree) getTallerChild(node *BinaryTreeNode) (child *BinaryTreeNode) {
	if node.left == nil {
		return node.right
	}
	if node.right == nil {
		return node.left
	}
	if node.left.height > node.right.height {
		return node.left
	} else {
		return node.right
	}
}

func (tree *BinaryTree) isBalance(node *BinaryTreeNode) bool {
	leftHeight := 0
	rightHeight := 0
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}
	return math.Abs(float64(leftHeight)-float64(rightHeight)) <= float64(1)
}

// 在删除后做的事情
func (tree *BinaryTree) afterRemove(node *BinaryTreeNode) {
	node = node.parent
	for node != nil {
		parent := node.parent
		if tree.isBalance(node) {
			tree.updateHeight(node)
		} else {
			tree.reBalance1(node)
		}
		node = parent
	}
}

// 在添加后做的事情
func (tree *BinaryTree) afterAdd(node *BinaryTreeNode) {
	if node == nil {
		return
	}
	node = node.parent
	for node != nil {
		parent := node.parent
		if tree.isBalance(node) {
			tree.updateHeight(node)
		} else {
			tree.reBalance1(node)
			break
		}
		node = parent
	}

}

// 红黑树
type RedBlackTree struct {
}
