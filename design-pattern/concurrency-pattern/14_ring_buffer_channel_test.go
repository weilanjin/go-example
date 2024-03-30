package concurrency_pattern

import (
	"log"
	"testing"
)

type ringBuffer struct {
	inCh  chan int
	outCh chan int
}

func NewRingBuffer(inCh, outCh chan int) *ringBuffer {
	return &ringBuffer{
		inCh:  inCh,
		outCh: outCh,
	}
}

func (r *ringBuffer) Run() {
	for v := range r.inCh {
		select {
		case r.outCh <- v:
		default:
			<-r.outCh
			r.outCh <- v
		}
	}
	close(r.outCh)
}

func Test14(t *testing.T) {
	inCh := make(chan int)
	outCh := make(chan int, 4)
	buffer := NewRingBuffer(inCh, outCh)
	go buffer.Run()
	for i := 0; i < 10; i++ {
		inCh <- i
	}
	close(inCh)
	for res := range outCh {
		log.Println(res)
	}
}
