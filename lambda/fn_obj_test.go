package lambda

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Fn interface {
	Calc(int) int
}

type Add struct {
}

func (a *Add) Calc(n int) int {
	return n + 1
}

func TestFnObject(t *testing.T) {
	calc := &Add{}
	data, err := json.Marshal(calc)
	if err != nil {
		t.Fatal(err)
	}
	_ = data

	var calc2 Fn
	if err = json.Unmarshal(data, calc2); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("calc2.Calc(2): %v\n", calc2.Calc(2))
}
