// Package semaphore
// 信号量(semaphore) 用来控制多个goroutine同时访问多个资源的同步原语
// 荷兰🇳🇱计算机科学家 Edsger Wybe Dijkstra 在 1963 年提出来的.
// - Dijkstra 算法是一种存在图形网络中找到最短路径的算法,被广泛用于网络路由和其他领域
// - “信号量” 同步机制,为并发提供了一种重要的工具
//
// 在操作系统中, 会给每一个进程分配一个信号量,代表每个进程目前的状态.
// 未得到控制权的进程,会在特定的地方被迫停下来,等待可以继续进行的型号到来.
// 信号量有两种类型:
// - 二元信号量 二元信号量只有两个值. 0 & 1 用于互斥访问共享资源(和互斥锁一样)
// - 计数信号量
package semaphore

/*
	P 操作(decrease、wait、acquire) 用来减少信号量的计数值
	V 操作(increase、signal、release) 增加信号量的计数值
	P passeren 通过、V vrijgeven 释放

    注:[]代表原子操作
    func V(semaphore S, integer I):
       [S <- S + I]

    func P(semaphore S, integer I):
       repeat:
          [if S >= I:
              S <- S - I
              break]

   n 个资源的池子
   P 操作相当于请求资源, 如果有足够的资源可用,就立即返回.
     如果没有资源或者资源不够,那么它可以不断尝试或者被阻塞等待.
   V 操作相当于释放资源, 把资源返回给信号.
*/