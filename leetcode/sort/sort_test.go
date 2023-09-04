package sort

import "testing"

var arr = []int{23, 32, 4, 17, 15, 28, 9, 5, 40, 41}

func TestBubbleSort(t *testing.T) {
	bubbleSort(arr)
	t.Log(arr)
}

func TestSelectionSort(t *testing.T) {
	selectionSort(arr)
	t.Log(arr)
}

func TestInsertionSort(t *testing.T) {
	insertionSort(arr)
	t.Log(arr)
}

func TestShellSort(t *testing.T) {
	shellSort(arr)
	t.Log(arr)
}

func TestMergeSort(t *testing.T) {
	arr = mergeSort(arr)
	t.Log(arr)
}

func TestQuickSort(t *testing.T) {
	quickSort(arr, 0, len(arr)-1)
	t.Log(arr)
}