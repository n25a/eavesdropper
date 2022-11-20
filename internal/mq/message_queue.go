package mq

type MessageQueue interface {
	Connect() error
	Close() error
	Publish(subject string, data interface{}) error
	Subscribe(subject string, handler func(payload interface{})) error
}

type MQType string

const (
	NATS MQType = "nats"
)
