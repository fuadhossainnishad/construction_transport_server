package rabbitmq

import (
    "construction_transport_server/internal/events"
    "context"
    "encoding/json"
    "log"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitEventPublisher struct {
    channel *amqp.Channel
    exchange string
}

func NewRabbitEventPublisher(conn *amqp.Connection) (*RabbitEventPublisher, error) {
    channel, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    exchange := "construction.events"
    err = channel.ExchangeDeclare(
        exchange,
        "topic",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return nil, err
    }
    return &RabbitEventPublisher{
        channel: channel,
        exchange: exchange,
    }, nil
}

func (p *RabbitEventPublisher) Publish(ctx context.Context, eventType events.EventType, payload interface{}) error {
    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    return p.channel.PublishWithContext(
        ctx,
        p.exchange,
        string(eventType),
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
            Timestamp:   time.Now(),
            Headers: amqp.Table{
                "event_type": string(eventType),
            },
        },
    )
}