package command

import "testing"

func TestCommand(t *testing.T) {
	mb := &MotherBoard{}
	startCommand := NewStartCommand(mb)
	rebootCommand := NewRebootCommand(mb)

	box1 := NewBox(startCommand, rebootCommand)
	box1.PressButton1()
	box1.PressButton2()

	box2 := NewBox(startCommand, rebootCommand)
	box2.PressButton1()
	box2.PressButton2()
}
