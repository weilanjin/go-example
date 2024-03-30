package factorymethod

import "testing"

func compute(factor OperatorFactory, a, b int) int {
	op := factor.Create()
	op.SetA(a)
	op.SetB(b)
	return op.Result()
}

func TestPlus(t *testing.T) {
	pof := PlusOperatorFactory{}
	if compute(pof, 12, 12) != 24 {
		t.Fatal("error with PlusOperatorFactory method pattern")
	}

	mof := MinusOperatorFactory{}
	if compute(mof, 12, 12) != 0 {
		t.Fatal("error with MinusOperatorFactory method pattern")
	}
}
