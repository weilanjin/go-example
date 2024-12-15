// Package syncmap
// map[k]v
// key 类型必须是可比较的 comparable, 可以通过 == 和 != 比较
// 可比较的数据类型: bool、整数、浮点数、字符串、指针、channel、接口、struct 和 数组都是可比较的
// 不可比较的数据类型: slice、map、function
// 主: 使用 struct 作为 key, struct 字段值容易被修改, 会导致 map 值容易找不到.
// 原生map:
// 1. 容易忘记被初始化
// 2. 无法保证线程安全
//
// 优化业务处理代码
// 1. 减少对锁的持有时间(锁的临界区 要尽可能的小)
// 2. 将串行的操作变成并行的子任务执行.
// -- 减少锁的粒度
// 1. 分片(shard), 将一个锁分成多个锁
// eg:
// https://github.com/orcaman/concurrent-map
// https://github.com/cornelk/hashmap
// https://github.com/alphadose/haxmap
package syncmap
