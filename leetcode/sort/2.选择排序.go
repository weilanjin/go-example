package sort

// 选择排序 - 逐渐缩小的区间，遍历每一轮找到最大(小)和端点位置交换
// 时间复杂度 最优 O(n^2), 最差 O(n^2)
// 空间复杂度 O(1)
func selectionSort(arr []int) {
	for i := len(arr) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if arr[j] > arr[i] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}