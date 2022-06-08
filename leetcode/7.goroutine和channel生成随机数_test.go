package leetcode

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

// 写代码实现两个 goroutine，其中⼀个产⽣随机数并写⼊到 go channel 中，另外⼀
//个从 channel 中读取数字并打印到标准输出。最终输出五个随机数。

func Test7(t *testing.T) {
	ch := make(chan int, 1)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch <- rand.Intn(5)
		}
		close(ch)
	}()
	go func() {
		defer wg.Done()
		for i := range ch {
			fmt.Print(i)
		}
		fmt.Println()
	}()
	wg.Wait()
}
