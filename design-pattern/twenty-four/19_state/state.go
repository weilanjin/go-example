package state

import "log"

// 状态模式 用于分离状态和行为
type Week interface {
	Today()
	Next(*DayContext)
}

type DayContext struct {
	today Week
}

func NewDayContext() *DayContext {
	return &DayContext{
		today: &Sunday{},
	}
}

func (d *DayContext) Today() {
	d.today.Today()
}

func (d *DayContext) Next() {
	d.today.Next(d)
}

type Sunday struct{}

func (*Sunday) Today() {
	log.Println("Sunday")
}

func (*Sunday) Next(context *DayContext) {
	context.today = &Monday{}
}

type Monday struct{}

func (*Monday) Today() {
	log.Println("Monday")
}

func (*Monday) Next(context *DayContext) {
	context.today = &Tuesday{}
}

type Tuesday struct{}

func (*Tuesday) Today() {
	log.Println("Tuesday")
}

func (*Tuesday) Next(context *DayContext) {
	context.today = &Wednesday{}
}

type Wednesday struct{}

func (*Wednesday) Today() {
	log.Println("Wednesday")
}

func (*Wednesday) Next(context *DayContext) {
	context.today = &Thursday{}
}

type Thursday struct{}

func (*Thursday) Today() {
	log.Println("Thursday")
}

func (*Thursday) Next(ctx *DayContext) {
	ctx.today = &Friday{}
}

type Friday struct{}

func (*Friday) Today() {
	log.Println("Friday")
}

func (*Friday) Next(ctx *DayContext) {
	ctx.today = &Saturday{}
}

type Saturday struct{}

func (*Saturday) Today() {
	log.Println("Saturday")
}

func (*Saturday) Next(ctx *DayContext) {
	ctx.today = &Sunday{}
}
