package synchronize

import (
	"log"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestMapGC(t *testing.T) {

	m := sync.Map{}
	for i := 0; i < 3000000; i++ {
		m.Store(i, "map_gc"+strconv.Itoa(i))
	}

	for i := 0; i < 2000000; i++ {
		m.Delete(i)
	}

	runtime.GC()
	log.Println("gc ------ 1")
	time.Sleep(time.Second * 8)
	log.Println("gc ------ 2")
	time.Sleep(time.Second * 8)
}
