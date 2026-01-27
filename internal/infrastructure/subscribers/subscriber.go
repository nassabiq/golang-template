package subscribers

import "github.com/nats-io/nats.go"

type Subscriber interface {
	Subject() string
	Durable() string
	Subscribe(js nats.JetStreamContext) error
}
