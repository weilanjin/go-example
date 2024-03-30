package memento

import "testing"

func TestMemento(t *testing.T) {
	g := &Game{
		hp: 10,
		mp: 10,
	}
	g.Status()
	progress := g.Save()

	g.Play(-2, -3)
	g.Status()

	g.Load(progress)
	g.Status()

	//18:06:11 Current HP:10. MP: 10
	//18:06:11 Current HP:7. MP: 8
	//18:06:11 Current HP:10. MP: 10
}
