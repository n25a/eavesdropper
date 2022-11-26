package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/n25a/eavesdropper/internal/app"

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

func (n *natsMQ) Subscribe(subject string, insertFunc insertFunction) error {
	sub, err := n.natsConnection.QueueSubscribe(subject, QGroup, natsHandler(insertFunc))
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

func natsHandler(insertFunc insertFunction) nats.MsgHandler {
	return func(msg *nats.Msg) {
		var data map[string]interface{}
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			log.Fatalln(err)
			return
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		var args []interface{}
		for _, field := range app.A.Schemas[msg.Subject].Fields {
			args = append(args, data[field])
		}

		err = insertFunc(ctx, app.A.Schemas[msg.Subject].Query, args...)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}
}
