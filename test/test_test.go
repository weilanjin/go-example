package test

import "testing"

func add(a, b int) int {
	return a + b
}

// go test -short
func TestShort(t *testing.T) {
	if testing.Short() {
		t.Skip("short 模式下跳过")
	}
	if sum := add(3, 4); sum != 7 {
		t.Errorf("expect %d, actual %d", 7, sum)
	}
}

// 单元测试
func TestUnit(t *testing.T) {
	var dataset = []struct {
		name string // 测试用例名称
		a    int
		b    int
		out  int
	}{
		{"case1", 1, 2, 3},
		{"case2", 2, 3, 5},
		{"case3", 3, 4, 7},
		{"case4", -9, 8, -1},
		{"case5", 0, 0, 0},
	}
	for _, v := range dataset {
		if sum := add(v.a, v.b); sum != v.out {
			t.Errorf("[%s]case %d: expect %d", v.name, sum, v.out)
		}
	}
}
