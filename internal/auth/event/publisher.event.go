package event

import (
	"construction_transport_server/infrastructure/messaging/rabbitmq"
	"context"
	"encoding/json"
	"log"
)

type EventPublisher interface {
	PublishUserRegistered(ctx context.Context, email, otp string)
	PublishBookingCreated(ctx context.Context, bookingID, customerID int64)
	PublishJobStatusChanged(ctx context.Context, bookingID int64, status string)
}

type eventPublisher struct {
	publisher *rabbitmq.Publisher
}

func NewEventPublisher(pub *rabbitmq.Publisher) EventPublisher {
	return &eventPublisher{publisher: pub}
}

func (e *eventPublisher) PublishUserRegistered(ctx context.Context, email, otp string) {
	body, _ := json.Marshal(map[string]string{
		"email": email,
		"otp":   otp,
	})
	go func() {
		err := e.publisher.Publish(ctx, "user.registered", body)
		if err != nil {
			log.Printf("failed to publish user.registered: %v", err)
		}
	}()
}

func (e *eventPublisher) PublishBookingCreated(ctx context.Context, bookingID, customerID int64) {
	body, _ := json.Marshal(map[string]interface{}{
		"booking_id":  bookingID,
		"customer_id": customerID,
	})
	go func() {
		err := e.publisher.Publish(ctx, "booking.created", body)
		if err != nil {
			log.Printf("failed to publish booking.created: %v", err)
		}
	}()
}

func (e *eventPublisher) PublishJobStatusChanged(ctx context.Context, bookingID int64, status string) {
	body, _ := json.Marshal(map[string]interface{}{
		"booking_id": bookingID,
		"status":     status,
	})
	go func() {
		err := e.publisher.Publish(ctx, "job.status.updated", body)
		if err != nil {
			log.Printf("failed to publish job.status.updated: %v", err)
		}
	}()
}
