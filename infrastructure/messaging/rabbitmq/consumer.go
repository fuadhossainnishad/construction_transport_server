package rabbitmq

import (
    "construction_transport_server/internal/events"
    "context"
    "encoding/json"
    "log"
    "sync"

    amqp "github.com/rabbitmq/amqp091-go"
)

type EventConsumer struct {
    conn     *amqp.Connection
    handlers map[events.EventType][]EventHandler
    mu       sync.RWMutex
}

type EventHandler func(ctx context.Context, payload []byte) error

func NewEventConsumer(conn *amqp.Connection) *EventConsumer {
    return &EventConsumer{
        conn:     conn,
        handlers: make(map[events.EventType][]EventHandler),
    }
}

func (c *EventConsumer) RegisterHandler(eventType events.EventType, handler EventHandler) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.handlers[eventType] = append(c.handlers[eventType], handler)
}

func (c *EventConsumer) Start(ctx context.Context) error {
    ch, err := c.conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()

    // Declare exchange
    err = ch.ExchangeDeclare("construction.events", "topic", true, false, false, false, nil)
    if err != nil {
        return err
    }

    // Create temporary queue for this consumer
    q, err := ch.QueueDeclare("", false, false, true, false, nil)
    if err != nil {
        return err
    }

    // Bind to all event types we have handlers for
    c.mu.RLock()
    for eventType := range c.handlers {
        err = ch.QueueBind(q.Name, string(eventType), "construction.events", false, nil)
        if err != nil {
            c.mu.RUnlock()
            return err
        }
    }
    c.mu.RUnlock()

    msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
    if err != nil {
        return err
    }

    go func() {
        for msg := range msgs {
            eventType := events.EventType(msg.RoutingKey)
            c.mu.RLock()
            handlers, ok := c.handlers[eventType]
            c.mu.RUnlock()
            if ok {
                for _, handler := range handlers {
                    go func(h EventHandler, data []byte) {
                        if err := h(ctx, data); err != nil {
                            log.Printf("Error handling event %s: %v", eventType, err)
                        }
                    }(handler, msg.Body)
                }
            }
        }
    }()
    return nil
}