package main

import (
	"log"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	mqURI       = "amqp://admin:admin@127.0.0.1:5672/my_vhost" // rabbitmq 环境变量设置的值-e RABBITMQ_DEFAULT_VHOST=my_vhost
	exchange    = "notification.ex"
	routingKey  = "compare_record.k"
	queueName   = "compare_record.q"
	consumerTag = "fusion"
)

func main() {
	r := gin.Default()
	r.GET("/ws/ntf/:id", notification)

	UseMQ() // init mq

	r.Run(":8080")
}

func UseMQ() {
	mq, err := NewMqReceive(mqURI, exchange, amqp.ExchangeDirect, queueName, routingKey, consumerTag)
	if err != nil {
		log.Printf("%+v", err)
	}
	go mq.subscribe(mq.queue.Name)
}
