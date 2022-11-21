package mq

import "github.com/nats-io/nats.go"

type MQType string

const (
	NATS MQType = "nats"
)

type MessageQueue interface {
	Connect() error
	Close() error
	Publish(subject string, data interface{}) error
	Subscribe(subject string, payload interface{}, handler func(payload interface{}) nats.MsgHandler) error
	UnSubscribe() error
}

func NewMessageQueue(mqType MQType) MessageQueue {
	if mqType == NATS {
		return NewNatsMQ()
	}
	return nil
}
