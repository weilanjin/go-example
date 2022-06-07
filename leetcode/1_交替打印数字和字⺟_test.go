package leetcode

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// 12AB34CD56EF78GH910IJ1112KL1314MN1516OP1718QR1920ST2122UV2324WX2526YZ2728
func Test1(t *testing.T) {
	number, letter := make(chan struct{}, 1), make(chan struct{}, 1)
	ch := make(chan string, 2)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		number <- struct{}{}
		for i := 1; i < 29; i = i + 2 {
			select {
			case <-number:
				ch <- strconv.Itoa(i)
				ch <- strconv.Itoa(i + 1)
				letter <- struct{}{}
			}
		}
		wg.Done()
	}()
	go func() {
		for i := 65; i < 91; i = i + 2 {
			select {
			case <-letter:
				ch <- string(rune(i))
				ch <- string(rune(i + 1))
				number <- struct{}{}
			}
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()
	var sb strings.Builder
	for s := range ch {
		sb.WriteString(s)
	}
	log.Println(sb.String())
}
