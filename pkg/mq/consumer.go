package mq

import (
	"fmt"
	"io"

	amqp "github.com/rabbitmq/amqp091-go"
)

const defaultPrefetchCount = 50

type Message struct {
	Body        []byte
	ContentType string
	tag         uint64
}

type Consumer interface {
	Receive() (*Message, error)
	Ack(msg *Message)
	Nack(msg *Message, requeue bool)
}

func NewRabbitConsumer(conn *amqp.Connection, queueName string) (Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error conn.Channel: %w", err)
	}

	// durable queue: match sender declaration
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error ch.QueueDeclare: %w", err)
	}

	// limit prefetch to prevent client-side message pileup on slow consumer
	if err = ch.Qos(defaultPrefetchCount, 0, false); err != nil {
		return nil, fmt.Errorf("error ch.Qos: %w", err)
	}

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error ch.Consume: %w", err)
	}
	return &RabbitConsumer{
		ch:   ch,
		msgs: msgs,
	}, nil
}

func (c *RabbitConsumer) Receive() (*Message, error) {
	msg, ok := <-c.msgs
	if !ok {
		return nil, io.EOF
	}
	return &Message{
		Body:        msg.Body,
		ContentType: msg.ContentType,
		tag:         msg.DeliveryTag,
	}, nil
}

func (c *RabbitConsumer) Ack(msg *Message) {
	if err := c.ch.Ack(msg.tag, true); err != nil { //nolint:staticcheck
		// TODO: log ack failed
	}
}

func (c *RabbitConsumer) Nack(msg *Message, requeue bool) {
	if err := c.ch.Nack(msg.tag, true, requeue); err != nil { //nolint:staticcheck
		// TODO: log nack failed
	}
}

type RabbitConsumer struct {
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
}
