package strategr

// 策略模式 定义一系列算法，让这些算法在运行时可以互换，使得分离算法
import "log"

type Payment struct {
	context  *PaymentContext
	strategy PaymentStrategy
}

type PaymentContext struct {
	Name, CardID string
	Money        int64
}

func NewPayment(name, cardId string, money int64, strategy PaymentStrategy) *Payment {
	return &Payment{
		context: &PaymentContext{
			Name:   name,
			CardID: cardId,
			Money:  money,
		},
		strategy: strategy,
	}
}

func (p *Payment) Pay() {
	p.strategy.Pay(p.context)
}

type PaymentStrategy interface {
	Pay(ctx *PaymentContext)
}

type Cash struct{}

func (*Cash) Pay(ctx *PaymentContext) {
	log.Printf("Pay %d to %s by cash", ctx.Money, ctx.Name)
}

type Bank struct{}

func (*Bank) Pay(ctx *PaymentContext) {
	log.Printf("Pay %d to %s by bank", ctx.Money, ctx.Name)
}
