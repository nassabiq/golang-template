package messaging

type EventBus interface {
	Publish(subject string, payload []byte) error
}
