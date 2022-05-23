package singleton

import (
	"sync"
	"testing"
)

func TestSingleton(t *testing.T) {
	s := GetInstance()
	s2 := GetInstance()
	if s != s2 {
		t.Fatal("instance is not equal")
	}
}

const parCount = 100

func TestParallelSingleton(t *testing.T) {
	c := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(parCount)
	s := [parCount]Singleton{}
	for i := 0; i < parCount; i++ {
		go func(i int) {
			<-c
			s[i] = GetInstance()
			wg.Done()
		}(i)
	}
	close(c)
	wg.Wait()
	for i := 1; i < parCount; i++ {
		if s[i] != s[i-1] {
			t.Fatal("instance is not equal")
		}
	}
}
