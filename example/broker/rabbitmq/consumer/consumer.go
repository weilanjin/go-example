package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	mqURI       = "amqp://admin:admin@127.0.0.1:5672/my_vhost" // rabbitmq 环境变量设置的值-e RABBITMQ_DEFAULT_VHOST=my_vhost
	exchange    = "notification.ex"
	routingKey  = "compare_record.k"
	queueName   = "compare_record.q"
	consumerTag = "fusion"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
)

func init() {
	var err error
	conn, err = amqp.Dial(mqURI)
	if err != nil {
		log.Fatalln(err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}

	if err = channel.ExchangeDeclare(
		exchange,            // name of the exchange
		amqp.ExchangeDirect, // type
		true,                // durable
		false,               // delete when complete
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		log.Fatalln(err)
	}

	queue, err := channel.QueueDeclare(
		queueName, // 队列临时名称
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err = channel.QueueBind(
		queue.Name, // name of the queue
		routingKey, // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	subscribe()
}

func subscribe() {
	deliveries, err := channel.Consume(
		queueName,   // name
		consumerTag, // consumerTag,
		false,       // autoAck
		false,       // exclusive
		false,       // noLocal
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalln(err)
	}

	for d := range deliveries {
		log.Printf("%s", d.Body)
	}
}
