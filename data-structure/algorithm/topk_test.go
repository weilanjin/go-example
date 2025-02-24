package algorithm

import (
	"container/heap"
	"fmt"
	"testing"
)

// 问题 ：给定一个长度为  的无序数组 nums ，请返回数组中最大的  个元素。

// min heap

type MinHeap []int

func (h MinHeap) Len() int { return len(h) }

func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }

func (h MinHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func TopK(nums []int, k int) []int {
	if k <= 0 {
		return nil
	}
	if k > len(nums) {
		return nums
	}
	h := &MinHeap{} // 创建一个最小堆
	heap.Init(h)
	for _, num := range nums {
		if h.Len() < k {
			heap.Push(h, num)
		} else if num > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, num)
		}
	}
	res := make([]int, k)
	for i := 0; i < k; i++ {
		res[k-1-i] = heap.Pop(h).(int)
	}
	return res
}

func TestTopK(t *testing.T) {
	nums := []int{3, 2, 1, 5, 6, 4}
	k := 2
	fmt.Println(TopK(nums, k)) // 输出: [6, 5]
}