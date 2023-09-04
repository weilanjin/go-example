package sort

// 希尔排序 - 按照折半分成多个小组（多轮宏观微调）
// 改进的插入排序算法
// 时间复杂度 最优 O(nlog2n), 最差 O(nlog2n)
func shellSort(arr []int) {
	gap := len(arr) / 2
	for gap > 0 {
		for i := gap; i < len(arr); i++ {
			p := arr[i]
			j := i - gap
			for j >= 0 && arr[j] > p {
				arr[j+gap] = arr[j]
				j -= gap
			}
			arr[j+gap] = p
		}
		gap /= 2
	}
}