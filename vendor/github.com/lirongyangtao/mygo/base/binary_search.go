package base

//二分搜索

func BinarySearch(data []int, value int) (index int, ok bool) {
	return binarySearch(data, 0, len(data), value)
}

func binarySearch(data []int, begin, end int, value int) (index int, ok bool) {
	for {
		if begin == end {
			return begin, false
		}
		mid := begin + (end-begin)/2
		if data[mid] == value {
			return mid, true
		} else if data[mid] < value {
			begin = mid + 1
		} else if data[mid] > value {
			end = mid

		}
	}
}

// 插入排序的二分查找接口实现
func binarySearchForInsertMerge1(lessSwap *lessSwap, begin, end int, findIndex int) (index int) {
	for begin < end {
		mid := begin + (end-begin)/2
		if lessSwap.Less(mid, findIndex) == E1LessE2 || lessSwap.Less(mid, findIndex) == E1EqualE2 {
			begin = mid + 1
		} else if lessSwap.Less(mid, findIndex) == E1GenerateE2 {
			end = mid
		}
	}
	return begin
}

// 插入排序优化的二分搜索，找到第一个大于value的值
func binarySearchForInsertMerge(data []int, begin, end int, value int) (index int) {
	for begin < end {
		mid := begin + (end-begin)/2
		if data[mid] <= value {
			begin = mid + 1
		} else if data[mid] > value {
			end = mid
		}
	}
	return begin
}
