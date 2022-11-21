package mq

import (
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/n25a/eavesdropper/internal/config"
)

const QGroup = "eavesdropper"

type natsMQ struct {
	natsConnection *nats.Conn
	subscriptions  []*nats.Subscription
}

func NewNatsMQ() MessageQueue {
	return &natsMQ{}
}

func (n *natsMQ) Connect() error {
	var err error
	// TODO: add Options
	n.natsConnection, err = nats.Connect(config.C.MQ.Conf.NatsAddress)
	return err
}

func (n *natsMQ) Close() error {
	n.natsConnection.Close()
	return nil
}

func (n *natsMQ) Publish(subject string, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return n.natsConnection.Publish(subject, dataBytes)
}

func (n *natsMQ) Subscribe(subject string, payload interface{},
	handler func(payload interface{}) nats.MsgHandler) error {
	sub, err := n.natsConnection.QueueSubscribe(subject, QGroup, handler(payload))
	if err != nil {
		return err
	}

	n.subscriptions = append(n.subscriptions, sub)
	return nil
}

func (n *natsMQ) UnSubscribe() error {
	for _, sub := range n.subscriptions {
		err := sub.Unsubscribe()
		if err != nil {
			return err
		}
	}
	return nil
}
