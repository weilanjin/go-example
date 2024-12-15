// Package cyclicbarrier 循环屏障
// https://github.com/marusama/cyclicbarrier
package cyclicbarrier

import (
	"context"
	"errors"
	"sync"
)

// 同步屏障 (Barrier) 对于一组goroutine,程序中的一个同步屏障意味着
// 任何goroutine执行到此后都必须等待,直到所有的goroutine都到达到此
// 点才可以继续执行下文.
//
// Barrier (屏障、栅栏) 拦截一组对象, 等对象齐了才打开它.
//
// CyclicBarrier 允许一组goroutine彼此等待,到达一个共同的检查点,
// 然后到达下一个同步点,循环使用.
//
/*
	CyclicBarrier 和 WaitGroup
	CyclicBarrier 更适合用在“数量固定的goroutine等待到达同一个检
	查点”的场景中,而且在放行goroutine之后,CyclicBarrier可以重复
	使用. WaitGroup 使用需要小心翼翼.

	WaitGroup 更适合用在“一个goroutine等待一组goroutine到达一个检查点”

	** CyclicBarrier 的参与者之间相互等待. 而WaitGroup一般都是父goroutine等待,干活的子goroutine之间不需要相互等待 **

	CyclicBarrier
		New(n)
	WaitGroup
		var wg WaitGroup
		wg.Add(n)

	CyclicBarrier
		Await
	WaitGroup
		wg.Done
		wg.Wait
		wg.Add(n) // 再调一次重用
*/

var (
	ErrBrokenBarrier = errors.New("broken barrier")
)

type CyclicBarrier interface {
	// Await 等待所有goroutine到达此点, 如果被ctx.Done()中断,则会返回ErrBrokenBarrier
	Await(ctx context.Context) error

	// Reset 重置循环屏障到初始状态,如果当前有等待者,那么它们会返回ErrBrokenBarrier
	Reset()

	// GetNumberWaiting 获取等待的数量
	GetNumberWaiting() int

	// GetParties 获取参与者的数量
	GetParties() int

	// IsBroken 循环屏障是否处于中断状态
	IsBroken() bool
}

type round struct {
	count    int           // 这一轮参与者goroutine的数量
	waitCh   chan struct{} // 这一轮等待channel
	brokeCh  chan struct{} // 广播用的channel
	isBroken bool          // 屏障是否被人为破坏
}

type cyclicBarrier struct {
	parties       int          // 参与者的数量
	barrierAction func() error // 屏障打开时需要调用的函数

	lock  sync.RWMutex
	round *round // 轮次
}

func New(parties int) CyclicBarrier {
	return NewWithAction(parties, nil)
}

func NewWithAction(parties int, action func() error) CyclicBarrier {
	if parties <= 0 {
		panic("parties must be positive number")
	}
	return &cyclicBarrier{
		parties: parties,
		round: &round{
			waitCh:  make(chan struct{}),
			brokeCh: make(chan struct{}),
		},
		barrierAction: action,
	}
}

func (b *cyclicBarrier) Await(ctx context.Context) error {
	var ctxDoneCh <-chan struct{}
	if ctx != nil {
		ctxDoneCh = ctx.Done()
	}
	// 检查ctx是否已经被取消或者超时
	select {
	case <-ctxDoneCh:
		return ctx.Err()
	default:
	}
	// 加锁
	b.lock.Lock()
	if b.round.isBroken { // 如果这一轮的等待和释放已经完成
		b.lock.Unlock()
		return ErrBrokenBarrier
	}
	// 在这一轮数据中将调用的参与者数量加1
	b.round.count++
	// 先保存这一轮的相关对象备用, 避免发生数据竞争,获取新一轮的对象
	waitCh := b.round.waitCh
	breakCh := b.round.brokeCh
	count := b.round.count

	b.lock.Unlock()
	//下面就不需要锁了, 因为已经获取到了这一轮的相关对象到了本地
	if count > b.parties { // 不能超过参与者数量
		panic("CyclicBarrier .Await is called more than count of parties")
	}
	// 如果当前的调用者不是最后一个调用者,则被阻塞等待
	if count < b.parties {
		// 等待发生下面的情况之一
		// 1. 最后一个调用者到来
		// 2. 人为的破坏了本轮的等待
		// 3. ctx 被完成
		select {
		case <-waitCh:
			return nil
		case <-ctxDoneCh:
			b.breakBarrier(true)
			return ctx.Err()
		case <-breakCh:
			return ErrBrokenBarrier
		}
	} else {
		// 如果当前的 goroutine 是最后一个调用者,则执行barrierAction函数(如果设置了)
		if b.barrierAction != nil {
			err := b.barrierAction()
			if err != nil {
				b.breakBarrier(true)
				return err
			}
		}
		// 重置屏障,因为它可循环使用, 重置之后可以继续使用,那就是下一轮的等待和释放
		b.reset(true)
		return nil
	}
}

func (b *cyclicBarrier) breakBarrier(needLock bool) {
	if needLock {
		b.lock.Lock()
		defer b.lock.Unlock()
	}
	if !b.round.isBroken {
		b.round.isBroken = true
		close(b.round.brokeCh)
	}
}

// safe = true 只需要把本轮的waitCh关闭即可
// unsafe = false 强制重置
func (b *cyclicBarrier) reset(safe bool) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if safe { // 广播,让等待的goroutine 继续执行
		close(b.round.waitCh)
	} else if b.round.count > 0 {
		b.breakBarrier(false)
	}
	// 创建新的一轮检查
	b.round = &round{
		waitCh:  make(chan struct{}),
		brokeCh: make(chan struct{}),
	}
}

func (b *cyclicBarrier) Reset() {
	b.reset(false)
}

func (b *cyclicBarrier) GetNumberWaiting() int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.round.count
}

func (b *cyclicBarrier) GetParties() int {
	return b.parties
}

func (b *cyclicBarrier) IsBroken() bool {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.round.isBroken
}