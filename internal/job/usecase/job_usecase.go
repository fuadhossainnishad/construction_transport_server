package usecase

import (
	"construction_transport_server/internal/booking/domain"
	"construction_transport_server/internal/booking/repository"
	"context"
	"errors"
	"time"
)

type JobUsecase struct {
	bookingRepo repository.BookingRepository
	wsHub       WebSocketHub
}

type WebSocketHub interface {
	SendToUser(userID int64, data interface{})
}

func NewJobUsecase(bookingRepo repository.BookingRepository, wsHub WebSocketHub) *JobUsecase {
	return &JobUsecase{bookingRepo: bookingRepo, wsHub: wsHub}
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
	// Update booking status
	if err := u.bookingRepo.UpdateStatus(ctx, bookingID, newStatus); err != nil {
		return err
	}
	// Add timeline entry
	if err := u.bookingRepo.AddStatusTimeline(ctx, bookingID, newStatus, notes); err != nil {
		// log but don't fail
	}
	// Real-time notification via WebSocket
	go u.wsHub.SendToUser(b.CustomerID, map[string]interface{}{
		"event":      "job_status_updated",
		"booking_id": bookingID,
		"status":     newStatus,
		"notes":      notes,
	})
	// If job completed, trigger payment and review availability
	if newStatus == string(domain.StatusCompleted) {
		go u.completeJob(ctx, bookingID)
	}
	return nil
}

func (u *JobUsecase) completeJob(ctx context.Context, bookingID int64) {
	// Mark completed_at timestamp
	updates := map[string]interface{}{
		"completed_at": time.Now(),
		"status":       domain.StatusCompleted,
	}
	_ = u.bookingRepo.Update(ctx, bookingID, updates)
	// TODO: trigger Stripe payout to transporter
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
