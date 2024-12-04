package semaphore

// 信号量实现
// - 初始化信号量: 设定资源的初始化数量
// - P 操作: 将信号量的计数值减k,如果新值为负数,那么调用者会阻塞并加入等待队列中;否则,调用者会继续执行,并且获取k个资源.
// - V 操作: 将信号量的计数值加k,如果先前的计数值为负,则说明有等待的P操作的调用者. V 操作会从等待队列中取出一个等待的调用者,唤醒它,让它继续执行.
//
// func runtime_Semacquire(s *uint32)
// func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)
// func runtime_Semrelease(s *uint32, handoff bool, skipframes int)
//
// 第三方库 https://github.com/marusama/semaphore 资源数量不是固定的,而是动态变化的

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(capacity int) *Semaphore {
	if capacity <= 0 {
		capacity = 1 // 容量为 1, 就变成了互斥锁
	}
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

// Acquire P 操作 请求一个资源
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release V 操作 释放资源
func (s *Semaphore) Release() {
	<-s.ch
}