package leetcode

import (
	"log"
	"sync"
	"testing"
	"time"
)

func Test11(t *testing.T) {
	wg := sync.WaitGroup{}
	c := make(chan struct{})
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int, done <-chan struct{}) {
			defer wg.Done()
			<-done
			log.Println(num)
		}(i, c)
	}
	if WaitTimeout(&wg, time.Second*5) {
		close(c)
		log.Println("timeout exit")
	}
	time.Sleep(time.Second * 10)
}

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	// 要求sync.WaitGroup支持timeout功能
	// 如果timeout到了超时时间返回true
	// 如果WaitGroup自然结束返回false
	ch := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return false
	case <-time.After(timeout):
		return true
	}
}
