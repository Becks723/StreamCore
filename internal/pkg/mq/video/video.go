package video

import (
	"context"
	"errors"

	"StreamCore/internal/pkg/mq/model"
	"StreamCore/pkg/mq"
	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

type VideoMQ interface {
	PublishVisitEvent(ctx context.Context, event *model.VisitEvent) error
	Consumer() (mq.Consumer, error)
}

func NewVideoMQ(conn *amqp.Connection) VideoMQ {
	queueName := "video_queue"
	sender, err := mq.NewRabbitSender(conn, queueName)
	if err != nil {
		sender = nil
		// TODO: log mq unavailable
	}
	return &videomq{
		sender:    sender,
		conn:      conn,
		queueName: queueName,
	}
}

var errSenderNotInitialized = errors.New("sender not initialized")

func (m *videomq) PublishVisitEvent(ctx context.Context, event *model.VisitEvent) error {
	if m.sender == nil {
		return errSenderNotInitialized
	}
	buffer, err := sonic.Marshal(event)
	if err != nil {
		return err
	}
	return m.sender.Send(ctx, buffer)
}

func (m *videomq) Consumer() (mq.Consumer, error) {
	return mq.NewRabbitConsumer(m.conn, m.queueName)
}

type videomq struct {
	sender    mq.Sender
	conn      *amqp.Connection
	queueName string
}
