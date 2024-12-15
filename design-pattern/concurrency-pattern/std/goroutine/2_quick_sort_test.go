package goroutine

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// go 并发并不一定最快
// 1. 如果是cpu敏感型的程序,则可以尝试使生成的goroutine的数量和cpu的逻辑核数相等.
// 2. 如果是io敏感型的程序,则可以尝试使生成的goroutine的数量和cpu的逻辑核数倍数或10倍.

// 快速排序中的分区,把a分成左右两部分,左边小于右边部分
// [1, 4, 7, 2, 8, 4, 9, 3, 5, 6, 2]
func partition(a []int, lo, hi int) int {
	pivot := a[hi] // 将最后一个值作为分界值
	slow := lo - 1
	for fast := lo; fast < hi; fast++ {
		if a[fast] < pivot { // 如果值小于分解值, 则挪到左边
			slow++
			a[fast], a[slow] = a[slow], a[fast] // fast slow 对应的值交换
		}
	}
	a[slow+1], a[hi] = a[hi], a[slow+1]
	return slow + 1
}

func quickSort(a []int, lo, hi int) {
	if lo >= hi {
		return
	}
	p := partition(a, lo, hi)
	quickSort(a, lo, p-1)
	quickSort(a, p+1, hi)
}

func quickSort_go(a []int, lo, hi int, done chan struct{}) {
	if lo >= hi {
		done <- struct{}{}
		return
	}
	p := partition(a, lo, hi)
	childDone := make(chan struct{}, 2)
	go quickSort_go(a, lo, p-1, childDone) // 启动一个 goroutine 快速排序左边
	go quickSort_go(a, p+1, hi, childDone) // 启动一个 goroutine 快速排序右边
	<-childDone
	<-childDone
	done <- struct{}{}
}

// depth 的最优配置, 可以接近cpu核数.
func quickSort_go_v2(a []int, lo, hi int, done chan struct{}, depth int) {
	if lo >= hi {
		done <- struct{}{}
		return
	}
	depth--
	p := partition(a, lo, hi)
	if depth > 0 {
		childDone := make(chan struct{}, 2)
		go quickSort_go_v2(a, lo, p-1, childDone, depth) // 启动一个 goroutine 快速排序左边
		go quickSort_go_v2(a, p+1, hi, childDone, depth) // 启动一个 goroutine 快速排序右边
		<-childDone
		<-childDone
	} else {
		quickSort(a, lo, p-1)
		quickSort(a, p+1, hi)
	}
	done <- struct{}{}
}

func TestQuickSort(t *testing.T) {
	arr := []int{1, 4, 7, 2, 8, 4, 9, 3, 5, 6, 2}
	quickSort(arr, 0, 10)
	fmt.Printf("arr: %v\n", arr)
}

// 1000万随机数排序, 3种快速排序对比
func TestBenchQuickSort(t *testing.T) {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var n int = 1e7 // 1kw

	testData1 := make([]int, 0, n)
	testData2 := make([]int, 0, n)
	testData3 := make([]int, 0, n)
	for i := 0; i < n; i++ {
		num := rd.Intn(n * 100)
		testData1 = append(testData1, num)
		testData2 = append(testData2, num)
		testData3 = append(testData3, num)
	}

	now := time.Now()
	quickSort(testData1, 0, n-1)
	fmt.Printf("串行 time.Since(now): %v\n", time.Since(now)) // 串行 time.Since(now): 1.498151542s

	now = time.Now()
	done := make(chan struct{}, 1)
	quickSort_go(testData2, 0, n-1, done)
	<-done
	fmt.Printf("并行 time.Since(now): %v\n", time.Since(now)) // 并行 time.Since(now): 2.846951833s

	now = time.Now()
	done = make(chan struct{}, 1)
	quickSort_go_v2(testData3, 0, n-1, done, 16)
	<-done
	fmt.Printf("并行,控制并发数量 time.Since(now): %v\n", time.Since(now)) // 并行,控制并发数量 time.Since(now): 261.425542ms
}
