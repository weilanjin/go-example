package facade

import "testing"

func TestAPI(t *testing.T) {
	expect := "A module running\nB module running"
	a := NewAPI()
	s := a.Test()
	if s != expect {
		t.Fatalf("expect %s, return %s", expect, s)
	}
}
