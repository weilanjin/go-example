// Package pool
// 对象池设计 (object pool pattern)
// 重复使用对象
// 1.可以从池中获取对象 Get
// 2.在不需要的时候还给池子 Put
//
// sync.Pool
// 池化的对象可能会被垃圾回收
// New: 创建对象
// Get: 获取对象
// Put: 返回对象
package pool
