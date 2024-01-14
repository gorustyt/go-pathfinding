package base

import (
	"math"
	"math/rand"
)

type sortI interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) int32
}

func checkAndGetSliceLen(s sortI) *lessSwap {
	return &lessSwap{
		Swap:   s.Swap,
		Less:   s.Less,
		Length: s.Len(),
	}
}

type lessSwap struct {
	Less   func(i, j int) int32
	Swap   func(i, j int) //交换
	Length int            //数组长度
}

// 基础类型
type SortInt []int

func (s *SortInt) Len() int {
	return len(*s)
}
func (s *SortInt) Swap(i, j int) {
	data := *s
	data[i], data[j] = data[j], data[i]

}
func (s *SortInt) Less(i, j int) int32 {
	data := *s
	if data[i] < data[j] {
		return E1LessE2
	} else if data[i] > data[j] {
		return E1GenerateE2
	} else {
		return E1EqualE2
	}
}

// ================================冒泡====================================
func BubbleSort(s sortI) {
	lessSwap := checkAndGetSliceLen(s)
	bubbleSort2(lessSwap)
}

func bubbleSort(lessSwap *lessSwap) {
	for i := 0; i < lessSwap.Length; i++ {
		for j := 1; j < lessSwap.Length-i; j++ {
			if lessSwap.Less(j-1, j) == E1GenerateE2 {
				lessSwap.Swap(j-1, j)
			}
		}
	}
}

// 优化：完全有序就终止
func bubbleSort1(lessSwap *lessSwap) {
	for i := 0; i < lessSwap.Length; i++ {
		sorted := true
		for j := 1; j < lessSwap.Length-i; j++ {
			if lessSwap.Less(j-1, j) == E1GenerateE2 {
				lessSwap.Swap(j-1, j)
				sorted = false
			}
		}
		if sorted {
			break
		}
	}
}

// 优化：后面排序就停止
func bubbleSort2(lessSwap *lessSwap) {
	for i := 0; i < lessSwap.Length; i++ {
		sortedIndex := 1
		end := lessSwap.Length - i
		for j := 1; j < end; j++ {
			if lessSwap.Less(j-1, j) == E1GenerateE2 {
				lessSwap.Swap(j-1, j)
				sortedIndex = j - 1
			}
		}
		end = sortedIndex
	}
}

// ================================选择====================================
func SelectSort(s sortI) {
	lessSwap := checkAndGetSliceLen(s)
	selectSort(lessSwap)
}

// 找出最大值和最后一个交换
func selectSort(lessSwap *lessSwap) {
	for i := 0; i < lessSwap.Length; i++ {
		last := lessSwap.Length - i - 1
		for j := 0; j < last; j++ {
			if lessSwap.Less(j, last) == E1GenerateE2 {
				lessSwap.Swap(j, last)
			}
		}
	}
}

// ================================插入====================================
func InsertSort(s sortI) {
	lessSwap := checkAndGetSliceLen(s)
	insertSort(lessSwap)
}

func insertSort(lessSwap *lessSwap) {
	for i := 1; i < lessSwap.Length; i++ {
		index := binarySearchForInsertMerge1(lessSwap, 0, i, i)
		lessSwap.Swap(index, i)
	}
}

func insertSort0(lessSwap *lessSwap) {
	for i := 1; i < lessSwap.Length; i++ {
		for j := i; j > 0; j-- {
			if lessSwap.Less(j, j-1) == E1GenerateE2 {
				lessSwap.Swap(j, j-1)
			}
		}
	}
}

func InsertSort1(data []int) {
	if len(data) < 2 {
		return
	}
	insertSort3(data)
}

// 插入排序通用实现
func insertSort2(data []int) {
	for i := 1; i < len(data); i++ {
		for j := i; j > 0; j-- {
			if data[j] < data[j-1] {
				data[j], data[j-1] = data[j-1], data[j]
			}
		}
	}
}

// 优化：i 往前挪动的时候进行二分搜索
func insertSort3(data []int) {
	for i := 1; i < len(data); i++ {
		index := binarySearchForInsertMerge(data, 0, i, data[i])
		tmp := data[i]
		for j := i; j > index; j-- {
			if data[j] < data[j-1] {
				data[j], data[j-1] = data[j-1], data[j]
			}
		}
		data[index] = tmp
	}
}

// ================================堆排序====================================
func HeapSort(s sortI) {
	lessSwap := checkAndGetSliceLen(s)
	heapSort(lessSwap, 0, lessSwap.Length)
}

// 找出最大值和最后一个交换
func heapSort(lessSwap *lessSwap, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) >> 1; i >= 0; i-- {
		binaryShiftDown(lessSwap, i, hi, first)
	}

	for i := hi - 1; i >= 0; i-- {
		lessSwap.Swap(first, first+i)
		binaryShiftDown(lessSwap, lo, i, first)
	}
}

// ================================归并排序 ====================================
func MergeSort1(data []int) {
	if len(data) < 2 {
		return
	}
	mergeSort2(data, 0, len(data))
}

// 自顶向下的merge
func mergeSort1(data []int, a, b int) {
	if b-a < 2 {
		return
	}
	mid := a + (b-a)/2
	mergeSort1(data, a, mid)
	mergeSort1(data, mid, b)
	merge1(data, a, mid, b)
}
func merge1(data []int, a, mid, b int) {
	arr := make([]int, b-a)
	pArr := 0
	pLeft := a
	pRight := mid
	for pLeft < mid && pRight < b {
		if data[pLeft] < data[pRight] {
			arr[pArr] = data[pLeft]
			pArr++
			pLeft++
		} else {
			arr[pArr] = data[pRight]
			pArr++
			pRight++
		}
	}
	//左侧小数组如果有剩余
	for pLeft < mid {
		arr[pArr] = data[pLeft]
		pArr++
		pLeft++
	}
	//右侧有数组如果有剩余
	for pRight < b {
		arr[pArr] = data[pRight]
		pArr++
		pRight++
	}
	//覆盖原来的数组
	for i := a; i < b; i++ {
		data[i] = arr[i-a]
	}
}

// 自底向上的merge
func mergeSort2(data []int, a, b int) {
	length := b - a
	for sz := 1; sz <= length; sz += sz {
		for i := 0; i+sz < length; i += sz + sz {
			merge1(data, i, i+sz-1, int(math.Min(float64(i+sz+sz-1), float64(length))))
		}
	}
}

// ================================快速排序====================================
func QuickSort(s sortI) {
	lessSwap := checkAndGetSliceLen(s)
	quickSort(lessSwap, 0, lessSwap.Length)
}
func quickSort(lessSwap *lessSwap, a, b int) {
	if b-a < 2 {
		return
	}
	piot := partition1(lessSwap, a, b)

	quickSort(lessSwap, a, piot)
	quickSort(lessSwap, piot+1, b)
}

// 寻找轴点
func choosePivot(lessSwap *lessSwap, begin, end int) int {
	length := end - begin
	return begin + rand.Intn(length)
}

// 更新过程中更新轴点
func partition(lessSwap *lessSwap, begin, end int) int {
	piot := choosePivot(lessSwap, begin, end)
	lessSwap.Swap(piot, begin)
	i := begin + 1
	j := end - 1
	for {
		for i <= j && lessSwap.Less(i, begin) == E1LessE2 {
			i++
		}
		for i <= j && lessSwap.Less(j, begin) == E1GenerateE2 {
			j--
		}
		if i > j {
			break
		}
		lessSwap.Swap(j, i)
		i++
		j--

	}
	lessSwap.Swap(j, begin)
	return j
}

// 优化1：如果有序，提前返回
func partition1(lessSwap *lessSwap, begin, end int) int {
	piot := choosePivot(lessSwap, begin, end)
	lessSwap.Swap(piot, begin)
	i := begin + 1
	j := end - 1
	//进行一遍，如果有序直接返回
	for i <= j && lessSwap.Less(i, begin) == E1LessE2 {
		i++
	}
	for i <= j && lessSwap.Less(j, begin) == E1GenerateE2 {
		j--
	}
	if i > j {
		lessSwap.Swap(j, begin)
		return j
	}
	lessSwap.Swap(j, i)
	i++
	j--
	for {
		for i <= j && lessSwap.Less(i, begin) == E1LessE2 {
			i++
		}
		for i <= j && lessSwap.Less(j, begin) == E1GenerateE2 {
			j--
		}
		if i > j {
			break
		}
		lessSwap.Swap(j, i)
		i++
		j--

	}
	lessSwap.Swap(j, begin)
	return j
}

// =====================快排通用实现======================
func QuickSort1(data []int) {
	if len(data) == 0 {
		return
	}
	quickSort1(data, 0, len(data))
}
func quickSort1(data []int, begin, end int) {
	if end-begin < 2 {
		return
	}
	piot := partition2(data, begin, end)
	quickSort1(data, begin, piot)
	quickSort1(data, piot+1, end)
}

// 标准快排,需要存储元素位置
func partition2(data []int, begin, end int) int {
	piot := rand.Intn(len(data))
	data[piot], data[begin] = data[piot], data[begin]
	tmp := data[begin]
	i := begin
	j := end - 1
	for i < j {
		for i < j && data[j] > tmp {
			j--
		}
		data[i] = data[j]
		i++
		for i < j && data[i] < tmp {
			i++
		}
		data[j] = data[i]
		j--
	}
	data[i] = tmp
	return i
}

// ================================希尔排序====================================
