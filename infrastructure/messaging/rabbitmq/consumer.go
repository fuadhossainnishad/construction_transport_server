package rabbitmq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type EmailEvent struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

// Consumer listens for user.registered queue and sends email (mock)
func StartEmailWorker(ctx context.Context, conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("user.registered", true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var event EmailEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("failed to unmarshal email event: %v", err)
				continue
			}
			// In real system, send email via SMTP or third‑party service.
			log.Printf("📧 Sending OTP %s to %s", event.OTP, event.Email)
		}
	}()
	return nil
}
