package question

import (
	"sync"
	"testing"
)

// 在携程执行 wg.Add
// 使用 sync.WaitGroup 副本
func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		go func(wg sync.WaitGroup, i int) {
			wg.Add(1)
			println(i)
			wg.Done()
		}(wg, i)
	}
	wg.Wait()
}
