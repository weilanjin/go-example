package adapter

import "testing"

func TestAdapter(t *testing.T) {
	a := NewAdaptee()
	t2 := NewAdapter(a)
	s := t2.Request()
	if s != "adaptee method" {
		t.Fatalf("expect 'adaptee method', return %s", s)
	}
}
