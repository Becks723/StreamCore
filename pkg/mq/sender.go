package mq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender interface {
	Send(ctx context.Context, buffer []byte) error
	Close()
}

func NewRabbitSender(conn *amqp.Connection, queueName string) (Sender, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &RabbitSender{
		conn:      conn,
		ch:        ch,
		queueName: queueName,
	}, nil
}

func (s *RabbitSender) Send(ctx context.Context, buffer []byte) error {
	err := s.ch.PublishWithContext(
		ctx,
		s.exchange,
		s.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        buffer,
		})
	if err != nil {
		return err
	}
	return nil
}

func (s *RabbitSender) Close() {
	s.conn.Close()
	s.ch.Close()
}

type RabbitSender struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
	exchange  string
}
