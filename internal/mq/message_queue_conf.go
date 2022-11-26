package mq

// Conf is the configuration for the message queue.
type Conf struct {
	NatsAddress string `json:"nats_address"`
}
