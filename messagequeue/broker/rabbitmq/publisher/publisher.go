package publisher

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	exchangeName    string
	bindingKey      string
	messageTypeName string
	amqpChan        *amqp.Channel
	amqpConn        *amqp.Connection
}

func NewPublisher(amqpConn *amqp.Connection) (Publisher, error) {
	ch, err := amqpConn.Channel()
	if err != nil {
		return Publisher{}, err
	}
	defer ch.Close()
	pub := Publisher{
		amqpConn:     amqpConn,
		amqpChan:     ch,
		exchangeName: "notification.ex",
		bindingKey:   "compare_record.k",
	}
	return pub, nil
}
