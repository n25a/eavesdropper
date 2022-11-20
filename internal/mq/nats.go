package mq

import (
	"github.com/nats-io/nats.go"

	"github.com/n25a/eavesdropper/internal/config"
)

type natsMQ struct {
	natsConnection *nats.Conn
}

func NewNatsMQ() MessageQueue {
	return &natsMQ{}
}

func (n *natsMQ) Connect() error {
	var err error
	n.natsConnection, err = nats.Connect(config.C.MQ.Conf.NatsAddress)
	return err
}

func (n *natsMQ) Close() error {
	n.natsConnection.Close()
	return nil
}

func (n *natsMQ) Publish(subject string, data interface{}) error {
	panic("Not implemented")
}

func (n *natsMQ) Subscribe(subject string, handler func(payload interface{})) error {
	panic("Not implemented")
}
