package prototype

import "testing"

type Type1 struct {
	name string
}

func (t *Type1) Clone() Cloneable {
	tc := *t
	return &tc
}

type Type2 struct {
	name string
}

func (t *Type2) Clone() Cloneable {
	tc := *t
	return &tc
}

var manager *PrototypeManager

func init() {
	manager = NewPrototypeManager()
	t1 := &Type1{
		name: "type1",
	}
	manager.Set("t1", t1)

	t2 := &Type2{
		name: "type2",
	}
	manager.Set("t2", t2)
}

func TestClone(t *testing.T) {
	c, ok := manager.Get("t1")
	if !ok {
		t.Fatal("error")
	}
	c2 := c.Clone()
	if c == c2 {
		t.Fatal("error! get clone not working")
	}
}

func TestCloneFromManager(t *testing.T) {
	c, ok := manager.Get("t2")
	if !ok {
		t.Fatal("error")
	}
	c2 := c.Clone()
	t2 := c2.(*Type2)
	if t2.name != "type2" {
		t.Fatal("error")
	}
}
