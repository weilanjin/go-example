package semaphore

import (
	"container/list"
	"context"
	"sync"
)

// go 扩展库中 使用 互斥锁🔒+ List 实现
// 在使用信号量时,最常见的几个错误
// - 请求了资源,但是忘记了释放它
// - 释放了从未请求的资源
// - 长时间持有一个资源(即使不需要它)
// - 不持有资源, 却直接使用它

type waiter struct {
	n     int64
	ready chan struct{}
}

type Weighted struct {
	size    int64      // 资源数量
	cur     int64      // 当前已使用的资源数量
	mu      sync.Mutex // P/V 操作时上锁
	waiters list.List  // waiter 列表
}

func NewWeighted(size int64) *Weighted {
	return &Weighted{size: size}
}

func (w *Weighted) Acquire(ctx context.Context, n int64) error {
	w.mu.Lock()
	// 快速路径: 如果有足够的资源, 则不考虑ctx.Done的状态,将cur加上n就返回
	if w.size-w.cur >= n && w.waiters.Len() == 0 {
		w.cur += n
		w.mu.Unlock()
		return nil
	}
	if n > w.size {
		w.mu.Unlock()
		// 依赖ctx的状态返回,否则一直等待
		<-ctx.Done()
		return ctx.Err()
	}
	// 否则,就需要把调用者加入等待队列中
	// 创建一个 ready chan, 以便通知环形
	ready := make(chan struct{})
	elem := w.waiters.PushBack(waiter{n: n, ready: ready})
	w.mu.Unlock()
	select {
	case <-ctx.Done():
		err := ctx.Err()
		w.mu.Lock()
		select {
		case <-ready: // 如果被唤醒了, 则忽略ctx的状态
			err = nil
		default: // 从 waiters 中移除自己
			isFront := w.waiters.Front() == elem
			w.waiters.Remove(elem)
			// 如果自己是队列中的第一个, 则看下一个waiter甚至更多的waiter需要的资源是否少,可以的得到满足
			if isFront && w.size > w.cur {
				w.notifyWaiters()
			}
		}
		w.mu.Unlock()
		return err
	case <-ready: // 被唤醒
		return nil
	}
}

func (w *Weighted) Release(n int64) {
	w.mu.Lock()
	w.cur -= n // 释放了N个资源
	if w.cur < 0 {
		w.mu.Unlock()
		panic("semaphore: released more than held")
	}
	w.notifyWaiters() // 唤醒 waiter
	w.mu.Unlock()
}

// TryAcquire 尝试获取资源,不会发生阻塞,所以也需要 Context
func (w *Weighted) TryAcquire(n int64) bool {
	w.mu.Lock()
	success := w.size-w.cur >= 0 && w.waiters.Len() == 0
	if success {
		w.cur += n
	}
	w.mu.Unlock()
	return success
}

func (w *Weighted) notifyWaiters() {
	for {
		next := w.waiters.Front()
		if next == nil {
			break // 没有 waiter 了
		}
		wa := next.Value.(waiter)
		if w.size-w.cur < wa.n {
			// 在没有充足的token提供给下一个waiter的情况下,没有继续查找, 而是停止
			// 主要是避免某个waiter饥饿
			break
		}
		w.cur += wa.n
		w.waiters.Remove(next)
		close(wa.ready)
	}
}