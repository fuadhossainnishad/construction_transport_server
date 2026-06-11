package event

import (
	"construction_transport_server/infrastructure/messaging/rabbitmq"
	"context"
)

type EventPublisherImpl struct {
	publisher *rabbitmq.Publisher
}

func NewEventPublisher(pub *rabbitmq.Publisher) EventPublisher {
	return &EventPublisherImpl{publisher: pub}
}

func (e *EventPublisherImpl) PublishUserRegistered(ctx context.Context, email string, otp string) {
	// Use goroutine to avoid blocking registration
	go func() {
		_ = e.publisher.PublishUserRegistered(ctx, email, otp)
	}()
}
