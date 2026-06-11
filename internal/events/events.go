package events

import "context"

type EventType string

const (
    UserRegisteredEvent   EventType = "user.registered"
    PasswordResetEvent    EventType = "password.reset"
    BookingCreatedEvent   EventType = "booking.created"
    JobStatusChangedEvent EventType = "job.status.changed"
    JobCompletedEvent     EventType = "job.completed"
)

type UserRegisteredPayload struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    OTP    string `json:"otp"`
}

type BookingCreatedPayload struct {
    BookingID     int64 `json:"booking_id"`
    CustomerID    int64 `json:"customer_id"`
    TransporterID int64 `json:"transporter_id"`
}

type JobStatusChangedPayload struct {
    BookingID int64  `json:"booking_id"`
    Status    string `json:"status"`
    Timestamp int64  `json:"timestamp"`
}

type JobCompletedPayload struct {
    BookingID int64   `json:"booking_id"`
    TotalPrice float64 `json:"total_price"`
    TransporterID int64 `json:"transporter_id"`
    CustomerID int64   `json:"customer_id"`
}

type EventPublisher interface {
    Publish(ctx context.Context, eventType EventType, payload interface{}) error
}