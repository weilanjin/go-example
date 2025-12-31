package main

const (
	noop   = iota
	add    // 加法 （arg1, arg2 int64）
	sub    // 减法 （arg1, arg2 int64）
	mul    // 乘法 （arg1, arg2 int64）
	div    // 除法 （arg1, arg2 int64）
	mod    // 取模/取余 （arg1, arg2 int64）
	neg    // 取反 （arg1 int64）
	inc    // 自增 （arg1 int64）
	dec    // 自减 （arg1 int64）
	lti    // 小于 （arg1, arg2 int64）
	lts    // 小于（arg1, arg2 []byte）
	eqi    // 等于 （arg1, arg2 int64）
	eqs    // 等于 （arg1, arg2 []byte）
	gti    // 大于 （arg1, arg2 int64）
	gts    // 大于 （arg1, arg2 []byte）
	not    // 取非 （arg1 int64）
	concat // 连接 （arg1, arg2 []byte)
	index  // 取字符（s []byte, idx int64）
	str    // 字符转字符串（ch int64）
	alloc  // 申请空间｜内存分配（size int64）
	read   // 读取端口数据 <port uint64> (arg1 []byte)
	write  // 写入端口数据 <port uint64> (arg1 []byte)
	pushi  // 入栈 （val int64）
	pushs  // 入栈 <len uint16>（val [len]byte)
	pusha  // 取参数入栈 <idx int16>
	jmp    // 跳转<delta int64>
	jz     // 如果为假跳转到<delta int64>
	call   // 调用函数<delta int64>
	ret    // 函数返回<narg uint16>
	halt   // 终止
)
