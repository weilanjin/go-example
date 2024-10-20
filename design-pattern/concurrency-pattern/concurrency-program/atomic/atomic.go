package atomic

/*
	AddXxxx 函数

	func AddInt32(addr *int32, delta int32) (new int32)
	func AddInt64(addr *int64, delta int64) (new int64)
	func AddUint32(addr *uint32, delta uint32) (new uint32)
	func AddUint64(addr *uint64, delta uint64) (new uint64)
	func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)

	delta 正数时, 代表增加, 负数时, 代表减少.
	无符号的整数
	AddUint32(&x, ^uint32(c-1)) // 利用了计算机原理中的补码原理, 变减法为加法
	AddUint32(&x, ^uint32(0)) // 减1

	CompareAndSwapXxxx 函数

	- 需要提供
		1. 操作地址
		2. 旧的值
		3. 新的值

		”判断为相等才替换“
		addr的指向值 和 old的指向的值 相等, 则替换为new, 返回 true
		否则, 返回 false

	func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
	func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
	func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
	func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
	func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)
	func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)

	SwapXxxx 函数

	- 粗暴的替换

	func SwapInt32(addr *int32, new int32) (old int32)
	func SwapInt64(addr *int64, new int64) (old int64)
	func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
	func SwapUint32(addr *uint32, new uint32) (old uint32)
	func SwapUint64(addr *uint64, new uint64) (old uint64)
	func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)

	LoadXxxx 函数

	- 取出 addr 地址中的值

	func LoadInt32(addr *int32) (val int32)
	func LoadInt64(addr *int64) (val int64)
	func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)
	func LoadUint32(addr *uint32) (val uint32)
	func LoadUint64(addr *uint64) (val uint64)
	func LoadUintptr(addr *uintptr) (val uintptr)

	StoreXxxx 函数

	- 把一个值存入指定的addr地址中

	func StoreInt32(addr *int32, new int32)
	func StoreInt64(addr *int64, new int64)
	func StorePointer(addr *unsafe.Pointer, new unsafe.Pointer)
	func StoreUint32(addr *uint32, new uint32)
	func StoreUint64(addr *uint64, new uint64)
	func StoreUintptr(addr *uintptr, new uintptr)

	- 原子地存取对象, 通常被用在配置变更等,对一个struct原子操作的场景

	Value
		func (v *Value) CompareAndSwap(old, new any) (swapped bool)
		func (v *Value) Load() (val any)
		func (v *Value) Store(val any)
		func (v *Value) Swap(new any) (old any)
*/

//  Bool、Int32、Int64、Uint、Uintptr、Uint64、Uintptr
