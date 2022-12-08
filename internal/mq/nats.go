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

const qGroup = "eavesdropper"

type natsMQ struct {
	natsConnection *nats.Conn
	subscriptions  []*nats.Subscription
}

// NewNatsMQ - create new nats message queue
func NewNatsMQ() MessageQueue {
	return &natsMQ{}
}

// Connect - connect to nats
func (n *natsMQ) Connect() error {
	var err error
	// TODO: add Options
	n.natsConnection, err = nats.Connect(config.C.MQ.Conf.NatsAddress)
	return err
}

// Close - close connection
func (n *natsMQ) Close() error {
	n.natsConnection.Close()
	return nil
}

// Publish - publish message to nats
func (n *natsMQ) Publish(subject string, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return n.natsConnection.Publish(subject, dataBytes)
}

// Subscribe - subscribe to nats
func (n *natsMQ) Subscribe(subject string, insertFunc insertFunction) error {
	sub, err := n.natsConnection.QueueSubscribe(subject, qGroup, natsHandler(insertFunc))
	if err != nil {
		return err
	}

	n.subscriptions = append(n.subscriptions, sub)
	return nil
}

// UnSubscribe - unsubscribe from nats
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

		for _, schema := range app.A.Schemas[msg.Subject] {
			var args []interface{}
			for _, field := range schema.Fields {
				args = append(args, data[field])
			}

			err = insertFunc(ctx, schema.Query, args...)
			if err != nil {
				log.Fatalln(err)
				return
			}
		}
	}
}
