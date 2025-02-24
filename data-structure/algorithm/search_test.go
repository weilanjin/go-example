package algorithm

// 二分查找
// arr 已排序
// 时间复杂度 O(logn)
func BinarySearch(arr []int, target int) int {
	i, j := 0, len(arr)-1 // 左闭右闭区间, i,j 分别指向数组首元素，尾元素
	for i <= j {
		mid := (i + j) / 2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			i = mid + 1
		} else {
			j = mid - 1
		}
	}
	return -1
}