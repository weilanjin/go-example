package waitgroup

import (
	"sync/atomic"
)

type noCopy struct{}

func runtime_Semrelease(s *uint32, handoff bool, skipframes int) {}
func runtime_Semacquire(s *uint32)                               {}

type WaitGroup struct {
	noCopy noCopy        // noCopy 是一个辅助字段, 主要用于辅助vet工具检查是否通过copy复制这个WaitGroup实例
	state  atomic.Uint64 // 高32位为计数器的值,低32位位waiter的数量; 记录计数器的值 + waiter的数量
	seam   uint32        // 信号量, 用来唤醒 waiter
}

func (wg *WaitGroup) Add(delta int) {
	state := wg.state.Add(uint64(delta) << 32) // 计数器的值加 delta 值
	v := int32(state >> 32)                    // 右移32位, 只保留计数器的值
	w := uint32(state)                         // waiter的数量
	if v < 0 {
		panic("sync: negative WaitGroup counter")
	}
	if w != 0 && delta > 0 && v == int32(delta) { // 有 waiter 还在等待的时候, 不应该再并发调用 Add
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	if v > 0 || w == 0 { // 成功, 返回
		return
	}
	if wg.state.Load() != state {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	wg.state.Store(0) // 计数器的值位0, 将waiter的计数清0, 并唤醒 waiter
	for ; w != 0; w-- {
		runtime_Semrelease(&wg.seam, false, 0)
	}
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	for {
		state := wg.state.Load()
		v := int32(state >> 32) // 得到计数器的值
		// w := uint32(state)      // 得到waiter的数量
		if v == 0 { // 计数器的值位0, 说明没有任务了, 直接返回
			return
		}
		// 增加waiter的数量
		if wg.state.CompareAndSwap(state, state+1) { // 把本goroutine加入waiter中, 加入成功, 被阻塞等待唤醒(因为可能同时有多个goroutine在调用Wait), 否则循环检查
			runtime_Semacquire(&wg.seam)
			if wg.state.Load() != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			return
		}
	}
}
