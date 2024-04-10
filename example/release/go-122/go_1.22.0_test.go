// release https://tip.golang.org/doc/go1.22
package release

import (
	"fmt"
	"testing"
)

// Fixing For Loops in Go 1.22 https://go.dev/blog/loopvar-preview
// 1. for 循环不再共享循环变量
// 2. 且支持整数范围
func TestFor(t *testing.T) {
	done := make(chan struct{})
	vals := []string{"a", "b", "c"}
	for _, v := range vals {
		go func() {
			fmt.Println(v) // c c c -> b a  sc
			done <- struct{}{}
		}()
	}
	for range vals {
		<-done
	}

	// for i := 0; i < n; i++ { ... }
	for i := range 10 {
		fmt.Println(10 - i) // 10 9 8 7 6 5 4 3 2 1
	}
	fmt.Println("go1.22 has lift-off!")
}