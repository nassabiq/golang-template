package nats

import "github.com/nats-io/nats.go"

type JetStreamBus struct {
	js nats.JetStreamContext
}

func NewJetStreamBus(nc *nats.Conn) (*JetStreamBus, error) {
	js, err := nc.JetStream()

	if err != nil {
		return nil, err
	}

	return &JetStreamBus{js: js}, nil
}

func (bus *JetStreamBus) Publish(subject string, payload []byte) error {
	_, err := bus.js.Publish(subject, payload)
	return err
}
