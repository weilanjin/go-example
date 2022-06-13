package leetcode

import (
	"sync"
	"time"
)

// GO⾥⾯MAP如何实现key不存在 get操作等待 直到key存在或者超时，保证并发安全.
type sp interface {
	// 如果该key读取的goroutine挂起，则唤醒
	// 此方法不会阻塞，时刻都可以立即执行并返回
	Out(key string, val any)

	//如果key不存在阻塞，等待key存在或者超时
	Rd(key string, timeout time.Duration) any
}

type item struct {
	ch    chan struct{}
	val   any
	isNew bool
}

type spMap struct {
	m map[string]*item
	sync.RWMutex
}

func (s *spMap) Out(key string, val any) {
	s.Lock()
	defer s.Unlock()
	if it, ok := s.m[key]; ok && it.isNew {
		it.isNew = false
		it.val = val
		it.ch <- struct{}{}
	} else {
		s.m[key] = &item{
			val:   val,
			isNew: false,
		}
	}
}

func (s *spMap) Rd(key string, timeout time.Duration) any {
	s.RLock()
	res, ok := s.m[key]
	s.RUnlock()
	if ok {
		return res.val
	}

	it := &item{
		ch:    make(chan struct{}),
		isNew: true,
	}

	s.Lock()
	s.m[key] = it
	s.Unlock()

	select {
	case <-it.ch:
		s.RLock()
		res := s.m[key]
		s.RUnlock()
		return res.val
	case <-time.After(timeout):
		return nil
	}
}
