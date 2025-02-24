package algorithm

import (
	"math"
	"sort"
)

// 选择排序
// 原理： 开启一个循环，每轮从未排序区间选择最小的元素，将其放到已排序区间的末尾
// 时间复杂度：O(n^2)
func SelectionSort(arr []int) {
	for i := 0; i < len(arr)-1; i++ { // 外循环 i—> n-1
		for j := i + 1; j < len(arr); j++ { // 内循环 找到未排序区间中的最小值
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

// 冒泡排序
// 原理： 通过连续地比较与交换相邻元素实现排序
// 时间复杂度：O(n^2)
func BubbleSort(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		exchange := false
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				exchange = true
			}
		}
		if !exchange { // 如果没有发生交换，则说明数组已经有序，可以提前结束排序
			return
		}
	}
}

// 插入排序
// 原理： 工作原理与手动整理一副牌的过程非常相似
// 时间复杂度：O(n^2)
func InsertionSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		for j := i; j > 0; j-- {
			if arr[j] < arr[j-1] {
				arr[j], arr[j-1] = arr[j-1], arr[j]
			}
		}
	}
}

// 快速排序
// 原理： 是一种基于分治策略的排序算法。
// 左子数组任意元素 <= 基准数 <= 右子数组任意元素
// 时间复杂度：O(nlogn)
// 空间复杂度：O(logn)
func QuickSort(arr []int) {
	if len(arr) <= 1 {
		return
	}
	quickSortHelper(arr, 0, len(arr)-1)
}

func quickSortHelper(arr []int, low, high int) {
	if low < high {
		// 获取分区点索引
		pivot := partition(arr, low, high)
		// 递归排序左半边
		quickSortHelper(arr, low, pivot-1)
		// 递归排序右半边
		quickSortHelper(arr, pivot+1, high)
	}
}

// 哨兵划分
func partition(arr []int, low, high int) int {
	// 选择中间元素作为基准（避免最坏情况）
	pivot := arr[(low+high)/2]
	// 将基准元素交换到末尾
	arr[(low+high)/2], arr[high] = arr[high], arr[(low+high)/2]
	i := low // 小于基准的元素的边界
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			// 将小于基准的元素交换到左边
			arr[j], arr[i] = arr[i], arr[j]
			i++
		}
	}

	// 将基准数放到正确的位置
	arr[i], arr[high] = arr[high], arr[i]
	return i
}

// 归并排序
// 原理： 将数组拆分为若干个子数组，然后递归地排序每个子数组，最后将排序后的子数组合并为整体
// 划分阶段、合并阶段
// 时间复杂度：O(nlogn)
// 空间复杂度：O(n)
func MergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	// 1.递归分割数组
	mid := len(arr) / 2
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])
	// 2.合并两个已排序的子数组
	return merge(left, right)
}

// 合并两个已排序的子数组
func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	// 比较两个子数组中的元素，并将较小的元素添加到结果中
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	// 将剩余的元素添加到结果中
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// 堆排序
// 原理： 堆排序是一种基于堆的数据结构排序算法。
// 时间复杂度：O(nlogn)
// 空间复杂度：O(1)
func HeapSort(arr []int) {
	n := len(arr)

	// 1.构建最大堆
	for i := n/2 - 1; i >= 0; i-- {
		heapify(arr, n, i)
	}
	// 2.依次提取堆顶元素（最大值），并调整堆
	for i := n - 1; i > 0; i-- {
		// 将堆顶元素（最大值）与最后一个元素交换
		arr[0], arr[i] = arr[i], arr[0]
		// 调整剩余元素，使其满足最大堆性质
		heapify(arr, i, 0)
	}
}

// 调整堆， 使其满足最大堆性质
func heapify(arr []int, n, i int) {
	largest := i // 假设当前节点为最大值
	left := 2*i + 1
	right := 2*i + 2
	// 如果左子节点比当前最大值大
	if left < n && arr[left] > arr[largest] {
		largest = left
	}
	if right < n && arr[right] > arr[largest] {
		largest = right
	}
	// 如果最大值不是当前节点，则交换并递归调整
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, n, largest)
	}
}

// 桶排序
// 原理： 桶排序是一种基于比较的排序算法。它通过设置一些具有大小顺序的痛，每个痛对应一个数据范围，将数据平均分配到各个桶中，然后对每个桶进行单独的排序。最后将各个桶中的数据合并到一起。
// 时间复杂度：O(n+k)
func BucketSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	// 1.找出最大值和最小值
	maxVal, minVal := arr[0], arr[0]
	for _, num := range arr {
		if num > maxVal {
			maxVal = num
		}
		if num < minVal {
			minVal = num
		}
	}
	// 如果所有的元素相同，则直接返回
	if maxVal == minVal {
		return arr
	}

	// 2.计算桶的数量和每个桶的范围
	bucketCount := int(math.Ceil(math.Sqrt(float64(len(arr))))) // 桶数量取数组长度的平方根
	bucketSize := int(math.Ceil(float64(maxVal-minVal+1) / float64(bucketCount)))

	// 3.初始化桶
	buckets := make([][]int, bucketCount)
	for i := range buckets {
		buckets[i] = make([]int, 0)
	}

	// 4.将元素分配到桶中
	for _, num := range arr {
		index := (num - minVal) / bucketSize
		// 如果元素超出桶的范围，则将其分配到最后一个桶
		if index >= bucketCount {
			index = bucketCount - 1
		}
		buckets[index] = append(buckets[index], num)
	}

	// 5.对每个桶进行单独的排序
	sorted := make([]int, 0, len(arr))
	for _, bucket := range buckets {
		sort.Ints(bucket)
		sorted = append(sorted, bucket...)
	}
	return sorted
}