package sort

// 冒泡排序 - 依次和相邻的下一个相比是否满足条件，交换位置
// |23|32|4|17|15|28|9|5|40|41|
// 时间复杂度 最优 O(n), 最差 O(n^2)
// 空间复杂度 O(1)
func bubbleSort(arr []int) {
	for i := len(arr) - 1; i > 0; i-- {
		exchange := false
		for j := 0; j < i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				exchange = true
			}
		}
		if !exchange {
			return
		}
	}
}