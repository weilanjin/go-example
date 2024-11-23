package channel

/*
	Don't communicate by sharing memory; share memory by communicating.
					- 不要通过共享内存来通信，而是通过通信来共享内存

	一个goroutine可以把数据的"所有权"交给另一个goroutine

	CSP Communicating Sequential Processes 通信顺序进程 - 描述并发系统中交互的一种模式

	几种场景

	1.「信息交流」: 把它当作并发的buffer或者队列, 解决 生产者(producer)-消费者问题(consumer)
	2.「数据传递」: 一个goroutine将数据交给另一个goroutine.相当于把数据的所有权交出去
	3.「信号通知」: 一个goroutine可以将信号(关闭中、已关闭、数据已准备好等)传递给另一个goroutine
	4.「任务编排」: 可以让一组goroutine按照一定的顺序并发或者串行地执行.
	5.「互斥锁」: 利用channel也可以实现互斥锁的机制.


+---------+-------+-----------+----------+---------+----------+--------+
|         | nil   | not empty | empty    |   full  | not full | closed |
+---------+-------+-----------+----------+---------+----------+--------+
| receive | 阻塞   | 读到值     | 阻塞     | 读到值   | 读到值     | 读完值  |
| send    | 阻塞   | 写入值     | 写入值    | 阻塞    | 写入值     | panic  |
| close   | panic | 正常关闭    | 正常关闭  | 正常关闭 | 正常关闭   | panic  |
+---------+-------+-----------+----------+---------+-----------+--------+

	channel 选择方法参考
	- 队共享资源的并发访问使用传统的同步原语
	- 复制的任务编排和消息传递使用channel
	- 消息通知机制使用channel,除非只想通知一个goroutine, 才使用Cond
	- 简单等待所有任务的完成使用WaitGroup, 也有channel的推崇者使用channel,它们都可以
	- 需要和select语句结合时,使用channel
	- 需要和超时配合时,使用channel和Context.
*/