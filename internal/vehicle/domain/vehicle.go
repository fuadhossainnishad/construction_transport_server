package domain

import "time"

type Vehicle struct {
    ID             int64     `json:"id"`
    TransporterID  int64     `json:"transporter_id"`
    TruckType      string    `json:"truck_type"`      // Dump Truck, Flatbed, Pickup, etc.
    LicensePlate   string    `json:"license_plate"`
    PayloadCap     string    `json:"payload_capacity"` // e.g., "10 tons"
    AxleCount      int       `json:"axle_count"`
    HourlyRate     float64   `json:"hourly_rate"`
    IsActive       bool      `json:"is_active"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type CreateVehicleInput struct {
    TruckType    string  `json:"truck_type" binding:"required"`
    LicensePlate string  `json:"license_plate" binding:"required"`
    PayloadCap   string  `json:"payload_capacity"`
    AxleCount    int     `json:"axle_count"`
    HourlyRate   float64 `json:"hourly_rate" binding:"required,gt=0"`
}

type UpdateVehicleInput struct {
    TruckType    *string  `json:"truck_type"`
    LicensePlate *string  `json:"license_plate"`
    PayloadCap   *string  `json:"payload_capacity"`
    AxleCount    *int     `json:"axle_count"`
    HourlyRate   *float64 `json:"hourly_rate"`
    IsActive     *bool    `json:"is_active"`
}