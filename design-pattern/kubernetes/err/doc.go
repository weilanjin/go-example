package err

// 错误处理的需求
// 1. 有时候对于报错的是有预期的, 要根据报错的原因不同进行不同的处理;
//    此时要方便判断出错误的原因, 此时错误处理也是业务的一部分
// 2. 错误要做记录, 并且要让看记录的人能够方便的了解错误发生的具体情况.
//    这时需要上下文信息以及调试链的信息.

// src/builtin/builtin.go
type error interface {
	Error() string
}

// 生成error
// 1. errors.New // 字符串不能大写, 不能以标点符号或换行符结尾
// 2. fmt.Errorf // 格式化字符串

// 错误包装(wrap)
// 为错误添加额外的信息, 如发生位置, 错误的原因等
// 错误包装后会生成新的错误
// 错误经过多次包装, 会形成错误链
// fmt.Errorf("%w", err) // 包装错误, 可以取出 (err.Unwrap()) 获取原始错误
// 如果是一层包装就是字符串, 如果是多层包装就是slice
// fmt.Errorf("%v", err) // 只是包含了其他错误信息, 不能解包

// 哨兵错误
// 在包中定义的错误就是哨兵错误, 是API的一部分.
// io.EOF、context.Canceled

// errors.Is
// error 经过包装后不能使用 == 判断, 要使用 errors.Is 判断
// errors.Is 只要错误链中的任何一个==就返回true
// 自定义错误不能 ==, 要实现 Is
