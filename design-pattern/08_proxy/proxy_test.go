package proxy

import (
	"testing"
)

func TestSubject(t *testing.T) {
	p := Proxy{}
	s := p.Do()
	if s != "pre:real:after" {
		t.Fatalf("expect 'pre:real:after'， return %s", s)
	}
}
