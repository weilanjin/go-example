package _map

import (
	"log"
	"sync"
	"testing"
)

func Test1(t *testing.T) {
	var m sync.Map
	m.LoadOrStore("a", 1)
	m.Delete("a")
	log.Println(m)
}
