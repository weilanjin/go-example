package algorithm

import "container/list"

/*
	给定三根柱子，记为 A、B 和 C 。起始状态下，柱子 A 上套着  个圆盘，它们从上到下按照从小到大的顺序排列。
	我们的任务是要把这  个圆盘移到柱子 C 上，并保持它们的原有顺序不变。

在移动圆盘的过程中，需要遵守以下规则。
	圆盘只能从一根柱子顶部拿出，从另一根柱子顶部放入。
	每次只能移动一个圆盘。
	小圆盘必须时刻位于大圆盘之上。
*/

func Hanoi(A, B, C *list.List, n int) {
	n = A.Len()
	// 将 A 顶部 n 个圆盘借助 B 移到 C
	dfsHanoi(n, A, B, C)
}

func move(src, dst *list.List) {
	val := src.Back() // 取出src最上层的圆盘
	dst.PushBack(val) // 将原盘放入 dst 顶部
	src.Remove(val)
}

func dfsHanoi(n int, src, buf, dst *list.List) {
	if n == 1 {
		move(src, dst)
		return
	}
	// 子问题 f(i-1): 将src顶部 i-1 个圆盘借助 dis 移到 dst
	dfsHanoi(n-1, src, dst, buf)
	// 子问题 f(i-1): 将src剩余一个圆盘移到 dst
	move(src, dst)
	// 子问题 f(i-1): 将buf顶部 i-1 个圆盘借助 src 移到 dst
	dfsHanoi(n-1, buf, src, dst)
}