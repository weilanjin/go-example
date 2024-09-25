package syncmap

import (
	"sync"
	"testing"
)

func TestSyncMap(t *testing.T) {
	m := sync.Map{}
	m.Load(1)
}
