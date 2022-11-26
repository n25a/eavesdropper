package mq

import (
	"context"
)

type insertFunction func(ctx context.Context, query string, arguments ...interface{}) error

type MQType string

const (
	NATS MQType = "nats"
)

type MessageQueue interface {
	Connect() error
	Close() error
	Publish(subject string, data interface{}) error
	Subscribe(subject string, insertFunc insertFunction) error
	UnSubscribe() error
}

func NewMessageQueue(mqType MQType) MessageQueue {
	if mqType == NATS {
		return NewNatsMQ()
	}
	return nil
}
