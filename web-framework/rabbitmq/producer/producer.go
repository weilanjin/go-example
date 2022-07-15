package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	mqURI      = "amqp://admin:admin@127.0.0.1:5672/my_vhost" // rabbitmq 环境变量设置的值-e RABBITMQ_DEFAULT_VHOST=my_vhost
	exchange   = "notification.ex"
	routingKey = "compare_record.k"
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

	if err := channel.ExchangeDeclare(
		exchange,            // name
		amqp.ExchangeDirect, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // noWait
		nil,                 // arguments
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
	t := time.NewTicker(time.Second * 5)
	rand.Seed(time.Now().Unix())
	for {
		randNum := rand.Intn(3) + 1
		msg := fmt.Sprintf(`{"id": "%s", "msg":"[%s] This is a message"}`, strconv.Itoa(randNum), time.Now().Local())
		publish(msg)
		<-t.C
	}
}

func publish(msg string) {
	log.Println(msg)
	if err := channel.Publish(
		exchange,   // 默认的 exchange
		routingKey, // routingKey
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	); err != nil {
		log.Println(err)
	}
}
