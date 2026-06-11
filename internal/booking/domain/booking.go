package domain

import "time"

type Booking struct {
	ID             int64
	CustomerID     int64
	TransporterID  *int64
	VehicleID      *int64
	Status         string
	PickupAddress  string
	PickupLat      *float64
	PickupLng      *float64
	DropoffAddress string
	DropoffLat     *float64
	DropoffLng     *float64
	ScheduledAt    *time.Time
	TotalPrice     *float64
	DistanceKm     *float64
	WorkNotes      string
	CreatedAt      time.Time
}

type BookingStatus string

const (
	StatusPending         BookingStatus = "pending"
	StatusAssigned        BookingStatus = "assigned"
	StatusHeadingToPickup BookingStatus = "heading_to_pickup"
	StatusArrivedPickup   BookingStatus = "arrived_at_pickup"
	StatusLoaded          BookingStatus = "loaded"
	StatusInTransit       BookingStatus = "in_transit"
	StatusDelivered       BookingStatus = "delivered"
	StatusCompleted       BookingStatus = "completed"
	StatusCancelled       BookingStatus = "cancelled"
)

type Booking struct {
	ID             int64      `json:"id"`
	CustomerID     int64      `json:"customer_id"`
	TransporterID  *int64     `json:"transporter_id,omitempty"`
	VehicleID      *int64     `json:"vehicle_id,omitempty"`
	Status         string     `json:"status"`
	PickupAddress  string     `json:"pickup_address"`
	PickupLat      *float64   `json:"pickup_lat,omitempty"`
	PickupLng      *float64   `json:"pickup_lng,omitempty"`
	DropoffAddress string     `json:"dropoff_address"`
	DropoffLat     *float64   `json:"dropoff_lat,omitempty"`
	DropoffLng     *float64   `json:"dropoff_lng,omitempty"`
	ScheduledAt    *time.Time `json:"scheduled_at,omitempty"`
	TotalPrice     *float64   `json:"total_price,omitempty"`
	DistanceKm     *float64   `json:"distance_km,omitempty"`
	WorkNotes      string     `json:"work_notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateBookingInput struct {
	VehicleID      int64    `json:"vehicle_id" binding:"required"`
	PickupAddress  string   `json:"pickup_address" binding:"required"`
	PickupLat      *float64 `json:"pickup_lat"`
	PickupLng      *float64 `json:"pickup_lng"`
	DropoffAddress string   `json:"dropoff_address" binding:"required"`
	DropoffLat     *float64 `json:"dropoff_lat"`
	DropoffLng     *float64 `json:"dropoff_lng"`
	ScheduledAt    *string  `json:"scheduled_at"` // RFC3339
	WorkNotes      string   `json:"work_notes"`
}
