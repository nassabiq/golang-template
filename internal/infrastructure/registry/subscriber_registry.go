package registry

import (
	"log"

	"github.com/nassabiq/golang-template/internal/infrastructure/subscribers"
	"github.com/nats-io/nats.go"
)

type Registry struct {
	subscribers []subscribers.Subscriber
}

func New() *Registry {
	return &Registry{}
}

func (r *Registry) Register(sub subscribers.Subscriber) {
	r.subscribers = append(r.subscribers, sub)
}

func (r *Registry) Run(js nats.JetStreamContext) {
	for _, sub := range r.subscribers {
		log.Println("🔔 subscribing:", sub.Subject())
		if err := sub.Subscribe(js); err != nil {
			log.Fatal(err)
		}
	}
}
