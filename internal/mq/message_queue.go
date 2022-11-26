package mq

import (
	"context"
)

type insertFunction func(ctx context.Context, query string, arguments ...interface{}) error

// MQType is the type of message queue.
type MQType string

// Message queue types
const (
	NATS MQType = "nats"
)

// MessageQueue is the interface that wraps the basic message queue operations.
type MessageQueue interface {
	Connect() error
	Close() error
	Publish(subject string, data interface{}) error
	Subscribe(subject string, insertFunc insertFunction) error
	UnSubscribe() error
}

// NewMessageQueue returns a new message queue instance.
func NewMessageQueue(mqType MQType) MessageQueue {
	if mqType == NATS {
		return NewNatsMQ()
	}
	return nil
}
