package concurrencyprogram

import (
	"fmt"
	"sync"
)

// fizz_buzz 问题
// 输入数字1到n, 满足下面条件
// 1.如果这个数字可以被3整除,则输出 "fizz"
// 2.如果这个数字可以被5整除,则输出 "buzz"
// 3.如果这个数字既被3整除,也被5整除,则输出 "fizzbuzz"
// 3.如果这个数字既不能被3整除,也不能被5整除,则输出这个数字本身

type FizzBuzz struct {
	n   int
	chs []chan int
	wg  sync.WaitGroup
}

func NewFizzBuzz(n int) *FizzBuzz {
	chs := make([]chan int, 4)
	for i := 0; i < 4; i++ {
		chs[i] = make(chan int, 1)
	}
	return &FizzBuzz{
		n:   n,
		chs: chs,
	}
}

func (f *FizzBuzz) start() {
	// f.wg.Add(4)
	// go f.fizz()
	// go f.buzz()
	// go f.fizzbuzz()
	// go f.number()
	f.wg.Wait()
}

// 只处理能被3整除的数,next <- v 表示交个下一个 goroutine 处理
func (f *FizzBuzz) fizz() {
	defer f.wg.Done()
	next := f.chs[1]
	for v := range f.chs[0] {
		if v%3 == 0 {
			fmt.Println("fizz")
		} else {
			next <- v
		}
	}
}
