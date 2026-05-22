package mq

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var errPublishNotConfirmed = errors.New("message publish not confirmed by broker")

type Sender interface {
	Send(ctx context.Context, buffer []byte) error
	Close()
}

func NewRabbitSender(conn *amqp.Connection, queueName string) (Sender, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// durable queue: survive broker restart
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// publisher confirm: ensure messages are accepted by broker
	if err = ch.Confirm(false); err != nil {
		return nil, err
	}
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	return &RabbitSender{
		conn:      conn,
		ch:        ch,
		queueName: queueName,
		confirms:  confirms,
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
			ContentType:  "text/plain",
			Body:         buffer,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		return err
	}

	// wait for confirmation
	select {
	case conf := <-s.confirms:
		if !conf.Ack {
			return errPublishNotConfirmed
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
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
	confirms  <-chan amqp.Confirmation
}
