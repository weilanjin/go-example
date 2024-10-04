package ctx

import "time"

type Context interface {
	// 返回这个 Context 被完成(done)的截止时间.
	// 如果没有截止时间, 则 ok 的值是false
	// 每次调用该方法时都会返回相同的结果
	Deadline() (deadline time.Time, ok bool)

	// Done 方法返回一个chanel对象.
	// 在 Context 被取消时, chanel 会被关闭(close).
	// 如果没有被撤销, Done 方法可能返回nil.
	// 多次调用都会返回相同的结果.
	Done() <-chan struct{}

	// 可以获取错误信息
	// 如果 Done 没有被关闭,那么Err将返回 nil
	// 否则, 返回一个非nil的error
	Err() error

	// Value 方法返回一个值, 该值与 key 关联.
	Value(key any) any
}

var (
	// 两个预定义对象,相当于一个空壳,一般用来做最初始的 Context 对象
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

type emptyCtx struct{}
