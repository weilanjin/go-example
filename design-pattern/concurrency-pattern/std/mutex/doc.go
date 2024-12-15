// Mutex 是一种互斥锁 Mutual Exclusion
// 是一种用于控制多线程对共享资源的竞争访问的同步机制.
//
// Lock 在编程语言中, 也将其称为锁 lock
// 当一个线程获得互斥锁时,它将阻止其他线程对该资源的访问,直到该线程释放互斥锁.
//
// synchronization primitive 同步原语
//
// 最早 1968 年 paper 《Solutions of a Problem in Concurrent Programming Control》中
// 首次实现了一种同步机制, 防止两个进程同时进入一个临界区(critical section) 后来被称为 Dijkstra 互斥算法.
//
// 竞争条件(race condition) 和 数据竞争(data race)
// - 竞争条件 (状态)
// 指的是在多线程环境中,由于操作顺序的不确定性,导致程序『执行结果不确定』的情况.
// - 数据竞争 (问题)
// 指的是在多线程环境中,由于操作顺序的不确定性,导致的「数据不一致」问题.
package mutex

// 同步原语 和 并发原语
/*
	- 同步原语 是一种用于控制多线程同时执行的操作. 实现并发操作, 如 并行计算 (e.g. 原子操作, 信号量, 互斥锁)
	- 并发原语 是一种用于控制多线程之间的执行顺序的操作, 实现同步操作, 如 线程间的数据传递 (e.g. 条件变量、消息队列、事件通知)
*/

// CAS compare and swap 原子操作
