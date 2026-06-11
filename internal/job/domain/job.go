package domain

import "time"

type Job struct {
    BookingID      int64     `json:"booking_id"`
    CustomerName   string    `json:"customer_name"`
    PickupAddress  string    `json:"pickup_address"`
    DropoffAddress string    `json:"dropoff_address"`
    Status         string    `json:"status"`
    ScheduledAt    *time.Time `json:"scheduled_at"`
    TotalPrice     *float64  `json:"total_price"`
    VehicleType    string    `json:"vehicle_type"`
    CreatedAt      time.Time `json:"created_at"`
}

type JobStatusUpdate struct {
    BookingID int64  `json:"booking_id"`
    Status    string `json:"status"`
    Notes     string `json:"notes,omitempty"`
}

type JobTimeline struct {
    Status    string    `json:"status"`
    Notes     string    `json:"notes"`
    CreatedAt time.Time `json:"created_at"`
}