package mq

import (
	"StreamCore/internal/pkg/mq/interaction"
	"StreamCore/internal/pkg/mq/video"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MQSet struct {
	Video       video.VideoMQ
	Interaction interaction.InteractionMQ
}

func NewMQSet(conn *amqp.Connection) *MQSet {
	return &MQSet{
		Video:       video.NewVideoMQ(conn),
		Interaction: interaction.NewInteractionMQ(conn),
	}
}
