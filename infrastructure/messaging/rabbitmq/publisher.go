package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp.Channel
}

func NewPublisher(ch *amqp.Channel) *Publisher {
	return &Publisher{
		ch: ch,
	}
}

type EmailEvent struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (p *Publisher) PublishUserRegistered(ctx context.Context, email string, otp string) error {
	event := EmailEvent{
		Email: email,
		OTP:   otp,
	}

	body, _ := json.Marshal(event)

	return p.ch.PublishWithContext(
		ctx,
		"",
		"user.registered",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
