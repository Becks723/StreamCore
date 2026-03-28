package mq

import (
	"fmt"
	"io"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Body        []byte
	ContentType string
	tag         uint64
}

type Consumer interface {
	Receive() (*Message, error)
	Ack(msg *Message)
}

func NewRabbitConsumer(conn *amqp.Connection, queueName string) (Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error conn.Channel: %w", err)
	}
	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error ch.QueueDeclare: %w", err)
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
		// TODO: log act failed, just log is ok
	}
}

type RabbitConsumer struct {
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
}
