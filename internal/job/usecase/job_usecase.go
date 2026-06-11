package usecase

import (
	"construction_transport_server/internal/booking/domain"
	"construction_transport_server/internal/booking/repository"
	"construction_transport_server/internal/events"
	"context"
	"errors"
	"time"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/transfer"
)

type JobUsecase struct {
	bookingRepo  repository.BookingRepository
	wsHub        WebSocketHub
	eventPub     events.EventPublisher
	stripeClient *stripe.Client
}

type WebSocketHub interface {
	SendToUser(userID int64, data interface{})
}

func NewJobUsecase(
	bookingRepo repository.BookingRepository,
	wsHub WebSocketHub,
	eventPub events.EventPublisher,
	stripeClient *stripe.Client,
) *JobUsecase {
	return &JobUsecase{
		bookingRepo:  bookingRepo,
		wsHub:        wsHub,
		eventPub:     eventPub,
		stripeClient: stripeClient,
	}
}

// GetJobDetails returns full job info including timeline
func (u *JobUsecase) GetJobDetails(ctx context.Context, bookingID, transporterID int64) (*domain.Booking, []domain.JobTimeline, error) {
	b, err := u.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || b == nil {
		return nil, nil, errors.New("job not found")
	}
	if b.TransporterID == nil || *b.TransporterID != transporterID {
		return nil, nil, errors.New("access denied")
	}
	timeline, err := u.bookingRepo.GetTimeline(ctx, bookingID)
	return b, timeline, err
}

func (u *JobUsecase) ListAssignedJobs(ctx context.Context, transporterID int64) ([]domain.Booking, error) {
	return u.bookingRepo.ListByTransporter(ctx, transporterID)
}

func (u *JobUsecase) UpdateJobStatus(ctx context.Context, bookingID, transporterID int64, newStatus, notes string) error {
	b, err := u.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || b == nil || b.TransporterID == nil || *b.TransporterID != transporterID {
		return errors.New("job not found or not assigned")
	}
	if !isValidTransition(b.Status, newStatus) {
		return errors.New("invalid status transition")
	}
	if err := u.bookingRepo.UpdateStatus(ctx, bookingID, newStatus); err != nil {
		return err
	}
	_ = u.bookingRepo.AddStatusTimeline(ctx, bookingID, newStatus, notes)

	// WebSocket notification to customer
	go u.wsHub.SendToUser(b.CustomerID, map[string]interface{}{
		"event":      "job_status_updated",
		"booking_id": bookingID,
		"status":     newStatus,
		"notes":      notes,
	})

	// Publish status changed event
	go u.eventPub.Publish(context.Background(), events.JobStatusChangedEvent, events.JobStatusChangedPayload{
		BookingID: bookingID,
		Status:    newStatus,
		Timestamp: time.Now().Unix(),
	})

	if newStatus == string(domain.StatusCompleted) {
		go u.completeJob(ctx, b)
	}
	return nil
}

func (u *JobUsecase) completeJob(ctx context.Context, booking *domain.Booking) {
	// Mark completed_at
	updates := map[string]interface{}{
		"completed_at": time.Now(),
		"status":       domain.StatusCompleted,
	}
	_ = u.bookingRepo.Update(ctx, booking.ID, updates)

	// Trigger Stripe transfer to transporter
	if booking.TotalPrice != nil && booking.TransporterID != nil {
		// In real scenario, fetch transporter's Stripe account ID from DB
		stripeAccountID := "acct_xxx"              // retrieve from auth table stripe_account_id
		amount := int64(*booking.TotalPrice * 100) // cents
		params := &stripe.TransferParams{
			Amount:      stripe.Int64(amount),
			Currency:    stripe.String("usd"),
			Destination: stripe.String(stripeAccountID),
		}
		_, err := transfer.New(params)
		if err != nil {
			// log error, maybe retry later
		}
	}

	// Publish job completed event for further processing (e.g., review, receipt)
	_ = u.eventPub.Publish(context.Background(), events.JobCompletedEvent, events.JobCompletedPayload{
		BookingID:     booking.ID,
		TransporterID: *booking.TransporterID,
		CustomerID:    booking.CustomerID,
		TotalPrice:    *booking.TotalPrice,
	})
}

// isValidTransition defines allowed status changes
func isValidTransition(current, new string) bool {
	transitions := map[string][]string{
		string(domain.StatusAssigned):        {string(domain.StatusHeadingToPickup), string(domain.StatusCancelled)},
		string(domain.StatusHeadingToPickup): {string(domain.StatusArrivedPickup)},
		string(domain.StatusArrivedPickup):   {string(domain.StatusLoaded)},
		string(domain.StatusLoaded):          {string(domain.StatusInTransit)},
		string(domain.StatusInTransit):       {string(domain.StatusDelivered)},
		string(domain.StatusDelivered):       {string(domain.StatusCompleted)},
	}
	allowed, ok := transitions[current]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == new {
			return true
		}
	}
	return false
}
