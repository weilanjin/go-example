package mem_test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"lovec.wlj/pkg/mem"
)

func TestMemConsumed(t *testing.T) {
	before := mem.MemConsumed()
	m := sync.Map{}
	for i := 0; i < 3000000; i++ {
		m.Store(i, "map_gc"+strconv.Itoa(i))
	}
	after := mem.MemConsumed()
	fmt.Printf("%.3fkb\n", float64(after-before)/1024)
}
