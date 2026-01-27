package event

import (
	"encoding/json"

	"github.com/nassabiq/golang-template/internal/infrastructure/messaging"
)

type Publisher struct {
	bus messaging.EventBus
}

func NewAuthPublisher(bus messaging.EventBus) *Publisher {
	return &Publisher{bus: bus}
}

func (publisher *Publisher) ForgotPassword(payload ForgotPasswordEvent) error {
	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	return publisher.bus.Publish(ForgotPasswordSubject, data)
}
