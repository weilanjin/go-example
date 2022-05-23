package builder

import "testing"

func TestBuilder1(t *testing.T) {
	b := &Builder1{}
	d := NewDirector(b)
	d.Construct()
	s := b.Result()
	if s != "123" {
		t.Fatalf("Builder1 fail expect 123, %s", s)
	}
}

func TestBuilder2(t *testing.T) {
	b := &Builder2{}
	d := NewDirector(b)
	d.Construct()
	i := b.Result()
	if i != 6 {
		t.Fatalf("Builder1 fail expect 6, %d", i)
	}
}
