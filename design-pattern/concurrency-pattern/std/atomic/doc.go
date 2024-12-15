//	写的地址基本上都是内存对齐的(aligned)
//	操作系统CPU和编译器
//	32bit write 的地址总是 4 的倍数
//	64bit write 的地址总是 8 的倍数
//
// 现代的多处理器多核系统, 由于缓存、指令重排
// 一个核更改地址值, 更新到主内存中之前, 是在多级缓存中存放.
// 其他核看到的数据可能不一样
//
// 多处理器多核系统使用内存屏障(memory fence 或 memory barrier)
// 一个写内存屏障会告诉处理器, 必须等到其管道中的未完成操作都被刷新到内存中.在进行其他操作,
// 此操作会使相关CPU缓存失效, 让它从主内存中刷新最新的值.
//
// 微软专家 Lockless Programming Considerations for Xbox 360 and Microsoft Windows
// atomic 操作的对象是一个地址, 需要把可寻址的变量的地址参数传递给函数, 而不是把变量的值传递给函数.
//
// func AddInt32(addr *int32, delta int32) (new int32)
// func AddInt64(addr *int64, delta int64) (new int64)
// func AddUint32(addr *uint32, delta uint32) (new uint32)
// func AddUint64(addr *uint64, delta uint64) (new uint64)
// func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)
//
// delta 正数时, 代表增加, 负数时, 代表减少.
// 无符号的整数
// AddUint32(&x, ^uint32(c-1)) // 利用了计算机原理中的补码原理, 变减法为加法
// AddUint32(&x, ^uint32(0)) // 减1
//
// uber-go/atomic
// Bool 提供了 Toggle() 方法 对布尔值类型的值做反转.
// Uint32 Sub() 减法、Inc() 加+1 和 Dec() 和减-1
// 提供了 Float32、Float64、Duration、String、Error 等类型
// 以上类型 提供 MarshalJSON() 和 UnmarshalJSON() 方法
package atomic
