package limit_test

import (
	"github.com/weilanjin/go-example/web-framework/limit"
	"log"
	"sync"
	"testing"
	"time"
)

func TestCunt(t *testing.T) {
	var wg sync.WaitGroup
	var lc limit.Counter
	lc.Set(3, time.Second) // 1s 内最多请求3次
	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.Println("创建请求：", i)
		go func(i int) {
			if lc.Allow() {
				log.Println("响应请求：", i)
			}
			wg.Done()
		}(i)
		time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()
}