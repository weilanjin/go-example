package question

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	type person struct {
		name string
	}

	var m map[person]int
	p := person{"lanjin.wei"}
	fmt.Println(m[p]) // 0 key 不存
}

func TestMapDelete(t *testing.T) {
	s := make(map[string]int)
	delete(s, "h")      // 删除不存在的 key 不会报错
	fmt.Println(s["h"]) // 0
}

// ----------------------------
func TestMapDelete1(t *testing.T) {
	var m = map[string]int{
		"A": 21,
		"B": 22,
		"C": 23,
	}
	counter := 0
	for k, v := range m {
		if counter == 0 {
			delete(m, "A")
		}
		counter++
		fmt.Println(k, v)
	}
	// map 遍历是无序的，如果第一次遍历到A那么 counter = 3，否则 counter = 2
	fmt.Println("counter is", counter) // 2 或 3
}