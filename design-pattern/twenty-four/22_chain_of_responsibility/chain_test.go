package chain

import "testing"

func TestChain(t *testing.T) {
	project := NewProjectManagerChain()
	dep := NewDepManagerChain()
	general := NewGeneralManagerChain()

	project.SetSuccessor(dep)
	dep.SetSuccessor(general)

	var m Manager = project

	m.HandleFeeRequest("bob", 400)
	m.HandleFeeRequest("tom", 1400)
	m.HandleFeeRequest("ada", 10000)
	m.HandleFeeRequest("flo", 400)
}
