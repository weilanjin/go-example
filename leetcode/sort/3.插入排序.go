package sort

// 插入排序 - 大于要插入的值都往后移
func insertionSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		p := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > p {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = p
	}
}