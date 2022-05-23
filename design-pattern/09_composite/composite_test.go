package composite

import "testing"

func TestComposite(t *testing.T) {
	c := NewComponent(CompositeNode, "root")
	c1 := NewComponent(CompositeNode, "c1")
	c2 := NewComponent(CompositeNode, "c2")
	c3 := NewComponent(CompositeNode, "c3")

	l1 := NewComponent(LeafNode, "l1")
	l2 := NewComponent(LeafNode, "l2")
	l3 := NewComponent(LeafNode, "l3")

	c.AddChild(c1)
	c.AddChild(c2)

	c1.AddChild(c3)
	c1.AddChild(l1)

	c2.AddChild(l2)
	c2.AddChild(l3)

	c.Print("")
}
