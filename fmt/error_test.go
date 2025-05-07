package fmt_test

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorf(t *testing.T) {
	wrapper := errors.New("inner error")

	err1 := fmt.Errorf("[Server] Error %s", "参数错误")
	err2 := fmt.Errorf("added context: %w", wrapper)
	fmt.Printf("err1=%v, err2=%v\n", err1, err2)
}