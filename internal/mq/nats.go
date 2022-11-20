package mq

type natsMQ struct{}

func NewNatsMQ() MessageQueue {
	return natsMQ{}
}

func (n *natsMQ) Connect() error {
	panic("Not implemented")
}

func (n *natsMQ) Close() error {
	panic("Not implemented")
}

func (n *natsMQ) Publish(subject string, data interface{}) error {
	panic("Not implemented")
}

func (n *natsMQ) Subscribe(subject string, handler func(payload interface{})) error {
	panic("Not implemented")
}
