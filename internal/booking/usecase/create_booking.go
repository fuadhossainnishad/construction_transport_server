package usecase

import (
    "construction_transport_server/internal/booking/domain"
    "construction_transport_server/internal/booking/repository"
    "construction_transport_server/internal/vehicle/repository" as vehicleRepo
    "context"
    "errors"
    "time"
)

type BookingUsecase struct {
    bookingRepo repository.BookingRepository
    vehicleRepo vehicleRepo.VehicleRepository
    eventPub    EventPublisher // same interface as before
}

type EventPublisher interface {
    PublishBookingCreated(ctx context.Context, bookingID int64, customerID int64)
}

func NewBookingUsecase(bookingRepo repository.BookingRepository, vehicleRepo vehicleRepo.VehicleRepository, eventPub EventPublisher) *BookingUsecase {
    return &BookingUsecase{
        bookingRepo: bookingRepo,
        vehicleRepo: vehicleRepo,
        eventPub:    eventPub,
    }
}

func (u *BookingUsecase) CreateBooking(ctx context.Context, customerID int64, input domain.CreateBookingInput) (*domain.Booking, error) {
    // verify vehicle exists and belongs to a transporter
    vehicle, err := u.vehicleRepo.GetByID(ctx, input.VehicleID)
    if err != nil || vehicle == nil {
        return nil, errors.New("invalid vehicle")
    }
    // For MVP: assign this vehicle's transporter automatically
    booking := &domain.Booking{
        CustomerID:     customerID,
        TransporterID:  &vehicle.TransporterID,
        VehicleID:      &vehicle.ID,
        Status:         string(domain.StatusPending),
        PickupAddress:  input.PickupAddress,
        PickupLat:      input.PickupLat,
        PickupLng:      input.PickupLng,
        DropoffAddress: input.DropoffAddress,
        DropoffLat:     input.DropoffLat,
        DropoffLng:     input.DropoffLng,
        WorkNotes:      input.WorkNotes,
    }
    if input.ScheduledAt != nil {
        t, err := time.Parse(time.RFC3339, *input.ScheduledAt)
        if err == nil {
            booking.ScheduledAt = &t
        }
    }
    // calculate distance & price (simplified placeholder)
    // actual: call external API for distance, then price = distance * vehicle.HourlyRate / 60
    // for now set dummy total
    dummyPrice := 47.0
    booking.TotalPrice = &dummyPrice

    err = u.bookingRepo.Create(ctx, booking)
    if err != nil {
        return nil, err
    }
    // publish event for potential real-time notification to transporter
    u.eventPub.PublishBookingCreated(ctx, booking.ID, customerID)
    return booking, nil
}

func (u *BookingUsecase) GetBooking(ctx context.Context, bookingID, userID int64, role string) (*domain.Booking, error) {
    b, err := u.bookingRepo.GetByID(ctx, bookingID)
    if err != nil || b == nil {
        return nil, errors.New("booking not found")
    }
    if role == "ADMIN" {
        return b, nil
    }
    if role == "USER" && b.CustomerID != userID {
        return nil, errors.New("access denied")
    }
    if role == "TRANSPORTER" && (b.TransporterID == nil || *b.TransporterID != userID) {
        return nil, errors.New("access denied")
    }
    return b, nil
}

func (u *BookingUsecase) ListCustomerBookings(ctx context.Context, customerID int64) ([]domain.Booking, error) {
    return u.bookingRepo.ListByCustomer(ctx, customerID)
}

func (u *BookingUsecase) ListTransporterBookings(ctx context.Context, transporterID int64) ([]domain.Booking, error) {
    return u.bookingRepo.ListByTransporter(ctx, transporterID)
}

func (u *BookingUsecase) CancelBooking(ctx context.Context, bookingID, customerID int64) error {
    b, err := u.bookingRepo.GetByID(ctx, bookingID)
    if err != nil || b == nil || b.CustomerID != customerID {
        return errors.New("booking not found or not yours")
    }
    if b.Status != string(domain.StatusPending) && b.Status != string(domain.StatusAssigned) {
        return errors.New("cannot cancel booking in current state")
    }
    return u.bookingRepo.UpdateStatus(ctx, bookingID, string(domain.StatusCancelled))
}