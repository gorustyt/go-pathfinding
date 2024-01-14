package base

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// 红黑树节点
type RBtreeNode struct {
	parent   *RBtreeNode //父节点
	left     *RBtreeNode //左子节点
	right    *RBtreeNode //右子节点
	element  any         //元素
	tmpIndex int         //在二叉树中索引,打印辅助
	IsRed    bool
}

func (node *RBtreeNode) IsLeaf() bool {
	return node.left == nil && node.right == nil
}

// 是否有2个孩子
func (node *RBtreeNode) HasTwoChildren() bool {
	if node == nil {
		return false
	}
	return node.left != nil && node.right != nil
}
func NewRBtreeNode(element interface{}, parent *RBtreeNode) *RBtreeNode {
	return &RBtreeNode{
		parent:  parent,
		element: element,
		IsRed:   true, //所有节点初始化默认为红色
	}
}

// 判断红黑树是红色
func (node *RBtreeNode) ColorIsRed() bool {
	if node == nil { //空节点代表黑色
		return false
	}
	return node.IsRed
}

// 判断红黑树是黑色
func (node *RBtreeNode) ColorIsBlack() bool {
	return !node.ColorIsRed()
}

// 给节点染色
func (node *RBtreeNode) SetColorRed() {
	if node == nil {
		return
	}
	node.IsRed = true
}

func (node *RBtreeNode) SetColorBlack() {
	if node == nil {
		return
	}
	node.IsRed = false
}

func (node *RBtreeNode) IsLeftChild() bool {
	return node.parent != nil && node == node.parent.left
}
func (node *RBtreeNode) IsRightChild() bool {
	return node.parent != nil && node == node.parent.right
}

func (node *RBtreeNode) Sibling() *RBtreeNode {
	if node.IsLeftChild() {
		return node.parent.right
	} else if node.IsRightChild() {
		return node.parent.left
	}
	return nil
}

type RbTree struct {
	size int
	Root *RBtreeNode
	cmp  CompareFunc //比较
}

func NewRbTree(cmp CompareFunc) *RbTree {
	return &RbTree{
		cmp: cmp,
	}
}
func (tree *RbTree) Add(ele any) {
	if tree.Root == nil {
		tree.Root = NewRBtreeNode(ele, nil)
		tree.size++
		return
	}
	node := tree.Root
	var parent *RBtreeNode
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
		newNode := NewRBtreeNode(ele, parent)
		parent.right = newNode
		tree.size++
		tree.afterAdd(newNode)
	} else if res == E1GenerateE2 {
		newNode := NewRBtreeNode(ele, parent)
		parent.left = newNode
		tree.size++
		tree.afterAdd(newNode)
	} else {
		parent.element = ele
	}

}
func (tree *RbTree) removeLeaf(node *RBtreeNode) {
	if node.parent == nil {
		tree.Root = nil
	} else if node == node.parent.left {
		node.parent.left = nil
	} else {
		node.parent.right = nil
	}
	tree.afterRemove(node, nil)
}

// 找后继者
func (tree *RbTree) FindSuccessor(node *RBtreeNode) *RBtreeNode {
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
func (tree *RbTree) FindNode(ele interface{}) (*RBtreeNode, bool) {
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
func (tree *RbTree) Remove(ele interface{}) {
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
		tree.afterRemove(node, replaceNode)
	} else if node.parent == nil {
		if node.left != nil {
			tree.Root = node.left
			node.left.parent = nil
		} else if node.right != nil {
			tree.Root = node.right
			node.right.parent = nil
		}
		tree.afterRemove(node, nil)
	} else { //度为1的节点
		tree.remove1(node)
	}
	tree.size--
}

// 移除度为1的节点
func (tree *RbTree) remove1(node *RBtreeNode) {
	var replaceNode *RBtreeNode
	if node == node.parent.left {
		if node.left != nil {
			node.parent.left = node.left
			node.left.parent = node.parent
			replaceNode = node.left
		} else {
			node.parent.left = node.right
			node.right.parent = node.parent
			replaceNode = node.right
		}
	} else {
		if node.right != nil {
			node.parent.right = node.right
			node.right.parent = node.parent
			replaceNode = node.right
		} else {
			node.parent.right = node.left
			node.left.parent = node.parent
			replaceNode = node.left
		}
	}
	tree.afterRemove(node, replaceNode)
}

func (tree *RbTree) afterRemove(node, replace *RBtreeNode) {
	if node.ColorIsRed() { //case1:删除节点是红色不要任何处理
		return
	}
	if replace.ColorIsRed() { //case2:取代节点是红色，直接染黑
		replace.SetColorBlack()
		return
	}
	if node.parent == nil { //case3:根节点
		return
	}
	parent := node.parent
	isLeft := parent.left == nil || node.IsLeftChild()
	sib := parent.left
	if sib == nil {
		sib = parent.right
	}
	if isLeft {
		if parent.Sibling().ColorIsRed() { //兄弟节点是红色
			parent.Sibling().SetColorBlack()
			parent.SetColorRed()
			tree.rotateLL(parent, node)
			sib = parent.right
		}
		//兄弟必然是黑的
		if sib.left.ColorIsBlack() && sib.right.ColorIsBlack() {
			parentIsBlack := parent.ColorIsBlack()
			parent.SetColorBlack()
			sib.SetColorRed()
			if parentIsBlack {
				tree.afterRemove(parent, nil)
			}
		} else { //兄弟必然有一个红色节点
			if sib.right.ColorIsBlack() {
				tree.rotateLL(sib, sib.right)
				sib = parent.right
			}
			if parent.ColorIsRed() {
				sib.SetColorRed()
			} else {
				sib.SetColorBlack()
			}
			parent.SetColorBlack()
			sib.right.SetColorBlack()
			tree.rotateRR(parent, node)
		}
	} else {
		if parent.Sibling().ColorIsRed() { //兄弟节点是红色
			parent.Sibling().SetColorBlack()
			parent.SetColorRed()
			tree.rotateRR(parent, node)
			sib = parent.left
		}
		//兄弟必然是黑的
		if sib.right.ColorIsBlack() && sib.left.ColorIsBlack() {
			parentIsBlack := parent.ColorIsBlack()
			parent.SetColorBlack()
			sib.SetColorRed()
			if parentIsBlack {
				tree.afterRemove(parent, nil)
			}
		} else { //兄弟必然有一个红色节点
			if sib.left.ColorIsBlack() {
				tree.rotateLL(sib, sib.left)
				sib = parent.left
			}
			if parent.ColorIsRed() {
				sib.SetColorRed()
			} else {
				sib.SetColorBlack()
			}
			parent.SetColorBlack()
			sib.left.SetColorBlack()
			tree.rotateRR(parent, node)
		}
	}
}

func (tree *RbTree) afterAdd(node *RBtreeNode) {
	parent := node.parent
	if parent == nil { //case1:根节点
		tree.Root.SetColorBlack()
		return
	}
	if !parent.ColorIsRed() { //case2:添加节点的父节点是黑色，直接返回
		return
	}
	grand := parent.parent
	uncle := parent.Sibling()
	//case3:uncle节点是红色
	if parent.Sibling().ColorIsRed() { //上溢
		parent.SetColorBlack()
		uncle.SetColorBlack()
		grand.SetColorRed()
		tree.afterAdd(grand)
	} else {
		if parent.IsLeftChild() {
			if node.IsLeftChild() { //LL
				parent.SetColorBlack()
				grand.SetColorRed()
				tree.rotateRR(grand, parent)
			} else if node.IsRightChild() { //LR
				grand.SetColorRed()
				node.SetColorBlack()
				tree.rotateLL(parent, node)
				tree.rotateRR(grand, parent)
			}
		} else {
			if node.IsLeftChild() { //RL
				grand.SetColorRed()
				node.SetColorBlack()
				tree.rotateRR(parent, node)
				tree.rotateLL(grand, parent)
			} else if node.IsRightChild() { //RR
				parent.SetColorBlack()
				grand.SetColorRed()
				tree.rotateLL(grand, parent)
			}
		}
	}
}

func (tree *RbTree) rotateLL(grand, parent *RBtreeNode) {
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

func (tree *RbTree) rotateRR(grand, parent *RBtreeNode) {
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

// 获取二叉树的每一层节点索引
func (tree *RbTree) getTreeNodeIndexMap() (nodeIndexMap map[int]map[int]interface{}, height int) {
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
		node = que.Poll().(*RBtreeNode)
		tmp, ok := nodeIndexMap[height]
		if !ok {
			tmp = map[int]interface{}{}
		}
		if node.ColorIsRed() {
			tmp[node.tmpIndex] = "\033[1;31m" + fmt.Sprintf("%v", node.element) + "\033[0m"
		} else {
			tmp[node.tmpIndex] = "\033[1;30m" + fmt.Sprintf("%v", node.element) + "\033[0m"
		}

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
func (tree *RbTree) TreePrint() {
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
