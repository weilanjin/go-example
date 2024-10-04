// Package ctx
// Go 最初提供了 golang.org/x/net/context包，但是被弃用了
// Go v1.7 正式吧Context加入标准库中
// Go v1.9 type alias 的新特性, 解决了 context 包的兼容性问题
//
// Context 使用场景
// - 上下文信息传递(request-scoped), 比如处理HTTP请求、在请求处理链路上传递信息.
// - 控制子 goroutine 的运行
// - 超时控制的方法调用
// - 可以撤销的方法调用
//
// Context 的时候, 有一些约定俗成的规则
// 1. 一般函数使用Context的时候, 会把这个参数放在第一个参数位置
// 2. 从来不把 nil 当作 Context 类型的参数值, 可以使用 context.Background()
// 创建一个空的上下文对象,但不要使用nil.
// 3. Context 只用来临时做函数之间的上下文透传, 不能持久化Context或者长久保存
// Context.
// 把Context持久化到数据库、本地文件、全局变量、缓存中都是错误用法.
// 4. key的类型不应该是字符串类型或者其他内建类型,否则在包之间使用Context的时候容易
// 产生冲突.使用WithValue时, key类型应该是自定义的.
// 5. 通常使用 struct{}作为底层类型来定义key的类型.
// exported key 的静态类型通常是接口或者指针,这样可以尽量减少内存分配
package ctx

// type Context = context.Context       // 定义别名
// type CancelFunc = context.CancelFunc // 定义别名
