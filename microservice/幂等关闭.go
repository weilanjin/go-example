package microservice

import "sync"

// https://antonz.org/idempotent-close/
// 幂等性
//
//	对一个对象重复调用某个操作时，不会导致变化或错误。
type Gate struct {
	mu     sync.Mutex // 确保线程安全
	closed bool
}

func NewGate() *Gate {
	return &Gate{}
}

func (g *Gate) Close() error {
	g.mu.Lock()
	if g.closed {
		g.mu.Unlock()
		return nil // 如果已经关闭，直接返回
	}

	g.closed = true
	g.mu.Unlock()

	// 这里可以添加关闭逻辑
	return nil
}

func ExampleGate_Close() {
	var g = NewGate()
	defer g.Close() // 确保在使用完毕后调用 Close 方法

	g.Close() // 可以多次调用 Close 方法而不会有副作用
}
