package trace

import (
	"sync"
)

var pool sync.Pool

func Get() *TraceLog {
	return pool.Get().(*TraceLog)
}

func Put(o *TraceLog) {
	o.reset()
	pool.Put(o)
}
