package rabbitmq

import (
	"errors"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

var ErrCannotConnectRabbitMQ = errors.New("cannot connect to rabbit")

func NewRabbitMQConn(mqURI string) (*amqp.Connection, error) {
	var count int64
	for {
		conn, err := amqp.Dial(mqURI)
		if err == nil {
			slog.Info("📫 connected to rabbitmq 🎉")
			return conn, nil
		}
		count++
		slog.Error("failed to connect to RabbitMq...", err, mqURI)
		if count > _retryTimes {
			return nil, ErrCannotConnectRabbitMQ
		}
		time.Sleep(_backOffSeconds * time.Second)
	}
}
