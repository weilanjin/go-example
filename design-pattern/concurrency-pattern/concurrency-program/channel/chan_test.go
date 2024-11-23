package channel

import "testing"

func TestSendNilChan(t *testing.T) {
	var ch chan int
	ch <- 1 // fatal error: all goroutines are asleep - deadlock!
	for v := range ch {
		t.Log(v)
	}
}