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
