package synchronize

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"lovec.wlj/pkg/mem"
)

func TestMapGC(t *testing.T) {
	m := sync.Map{}
	begin := mem.MemConsumed()
	for i := 0; i < 3000000; i++ {
		m.Store(i, "map_gc"+strconv.Itoa(i))
	}
	before := mem.MemConsumed()
	for i := 0; i < 2000000; i++ {
		m.Delete(i)
	}
	after := mem.MemConsumed()
	fmt.Printf("begin memory %.3fkb\n", float64(begin)/1024)
	fmt.Printf("store() memory %.3fkb\n", float64(before)/1024)
	fmt.Printf("delete() memory %.3fkb\n", float64(after)/1024)
}
