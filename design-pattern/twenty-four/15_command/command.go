package command

import "log"

// 命令模式 把某个对象的方法调用封装到对象中，方便传递、存储、调用。

type Command interface {
	Execute()
}

type StartCommand struct {
	mb *MotherBoard
}

func NewStartCommand(mb *MotherBoard) *StartCommand {
	return &StartCommand{
		mb: mb,
	}
}

func (c *StartCommand) Execute() {
	c.mb.Start()
}

type RebootCommand struct {
	mb *MotherBoard
}

func NewRebootCommand(mb *MotherBoard) *RebootCommand {
	return &RebootCommand{
		mb: mb,
	}
}

func (c *RebootCommand) Execute() {
	c.mb.Reboot()
}

type MotherBoard struct{}

func (*MotherBoard) Start() {
	log.Println("system, starting")
}

func (*MotherBoard) Reboot() {
	log.Println("system rebooting")
}

type Box struct {
	button1 Command
	button2 Command
}

func NewBox(button1, button2 Command) *Box {
	return &Box{
		button1: button1,
		button2: button2,
	}
}

func (b Box) PressButton1() {
	b.button1.Execute()
}

func (b Box) PressButton2() {
	b.button2.Execute()
}
