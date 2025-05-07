package pprof_test

import (
	"log"
	"runtime/pprof"
	"testing"
	"time"
)

func TestGoroutine(t *testing.T) {
	go func() {
		goroutines := pprof.Lookup("goroutine")
		for range time.Tick(1 * time.Second) {
			// 每秒统计，当前goroutine数
			log.Printf("goroutine count: %d\n", goroutines.Count())
		}
	}()
	var blockForever chan struct{}
	for i := 0; i < 10; i++ {
		go func() { <-blockForever }()
		time.Sleep(500 * time.Millisecond)
	}
}
