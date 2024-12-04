package cyclicbarrier

import (
	"context"
	"github.com/marusama/cyclicbarrier"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCyclicBarrier(t *testing.T) {
	cnt := 0
	b := cyclicbarrier.NewWithAction(10, func() error {
		// 屏障打开时计数器的值+1
		cnt++
		return nil
	})
	var wg sync.WaitGroup
	wg.Add(10)
	for i := range 10 {
		go func(r int) {
			defer wg.Done()
			for j := range 5 {
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
				log.Printf("goroutine %d 来到第 %d 轮屏障", r, j)

				if err := b.Await(context.TODO()); err != nil {
					log.Fatalf("goroutine %d 遇到错误: %s", r, err)
				}
				log.Printf("goroutine %d 离开第 %d 轮屏障", r, j)
			}
		}(i)
	}

	wg.Wait()
	log.Printf("计数器的值: %d", cnt)
}