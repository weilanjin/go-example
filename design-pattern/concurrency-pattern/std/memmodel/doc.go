// Package memmodel
// https://research.swtch.com/gomm
// GO 内存模型
// 「并发环境中多个goroutine读取相同变量时,对变量可见性的保证」
// 也就是说什么条件下,一个goroutine在读取一个变量值时,能够看到其他goroutine对这个变量进行写的结果
// CPU内存指令重排、多级缓存的存在,保证多线程访问一个变量非常复杂
// 不同的CPU架构(x86、amd64、arm、power等)处理方式不一样
package memmodel

/*
	If you must read the rest of this document to understand the behavior of your program. you are being too clever
	Don't be clever.

	如果你必须阅读此文档的其余部分才能理解程序的行为. 那么您过于聪明了.
	不要太聪明.
*/

// sequenced before、synchronized before、happens before