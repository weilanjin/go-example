package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MqReceive struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	tag     string
	done    chan error
}

func NewMqReceive(mqURI, exchange, exchangeType, queueName, routingKey, consumerTag string) (*MqReceive, error) {
	mq := &MqReceive{
		tag:  consumerTag,
		done: make(chan error),
	}

	var err error
	if mq.conn, err = amqp.Dial(mqURI); err != nil {
		return nil, errors.WithStack(err)
	}

	if mq.channel, err = mq.conn.Channel(); err != nil {
		return nil, errors.WithStack(err)
	}

	// 绑定 exchange
	if err = mq.channel.ExchangeDeclare(
		exchange,            // name of the exchange
		amqp.ExchangeDirect, // type
		true,                // durable
		false,               // delete when complete
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		return nil, errors.WithStack(err)
	}

	// 绑定 queue
	mq.queue, err = mq.channel.QueueDeclare(
		queueName, // 队列临时名称
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// queue和exchange绑定
	if err = mq.channel.QueueBind(
		mq.queue.Name, // name of the queue
		routingKey,    // bindingKey
		exchange,      // sourceExchange
		false,         // noWait
		nil,           // arguments
	); err != nil {
		return nil, errors.WithStack(err)
	}
	return mq, nil
}

func (mq *MqReceive) subscribe(queue string) {
	deliveries, err := mq.channel.Consume(
		queue,  // name
		mq.tag, // consumerTag,
		false,  // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // arguments
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%+v", errors.New("Listen .... queue : "+queue))
	for d := range deliveries {
		var record *Record
		if err := json.Unmarshal(d.Body, &record); err != nil {
			log.Printf("%+v, %v", errors.New(string(d.Body)), err)
			continue
		}

		mq.WsSend(record.Id, string(d.Body))
	}
}

// id == websocketId
type Record struct {
	Id, Msg string
}

// func (mq *MqReceive) dataUnpacket(data []byte) (record *Record) {
// 	err := json.Unmarshal(data, &record)
// 	if err != nil {
// 		log.Printf("%+v", errors.New(websocketId+" : not's exist"))
// 		return
// 	}

// }

// 向websocket发送消息
func (mq *MqReceive) WsSend(websocketId, msg string) {
	_ws, ok := WsConns.Load(websocketId)
	if !ok {
		log.Printf("%+v", errors.New(websocketId+" : not's exist"))
		return
	}

	ws := _ws.(*websocket.Conn)

	if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Printf("%+v", errors.WithStack(err))
		WsConns.Delete(websocketId) // TODO 是否是中断连接异常
		return
	}
}
