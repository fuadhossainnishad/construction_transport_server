package usecase

import (
    "construction_transport_server/internal/booking/domain"
    "construction_transport_server/internal/booking/repository"
    "context"
    "errors"
)

type JobUsecase struct {
    bookingRepo repository.BookingRepository
    wsHub       WebSocketHub // for real-time updates
}

type WebSocketHub interface {
    SendToUser(userID int64, data interface{})
}

func NewJobUsecase(bookingRepo repository.BookingRepository, wsHub WebSocketHub) *JobUsecase {
    return &JobUsecase{bookingRepo: bookingRepo, wsHub: wsHub}
}

func (u *JobUsecase) ListAssignedJobs(ctx context.Context, transporterID int64) ([]domain.Booking, error) {
    return u.bookingRepo.ListByTransporter(ctx, transporterID)
}

func (u *JobUsecase) UpdateJobStatus(ctx context.Context, bookingID, transporterID int64, newStatus string) error {
    b, err := u.bookingRepo.GetByID(ctx, bookingID)
    if err != nil || b == nil || b.TransporterID == nil || *b.TransporterID != transporterID {
        return errors.New("job not found or not assigned to you")
    }
    // validate allowed status transitions
    if !isValidTransition(b.Status, newStatus) {
        return errors.New("invalid status transition")
    }
    err = u.bookingRepo.UpdateStatus(ctx, bookingID, newStatus)
    if err != nil {
        return err
    }
    // insert into status timeline
    go u.insertStatusTimeline(bookingID, newStatus)

    // notify customer via WebSocket
    u.wsHub.SendToUser(b.CustomerID, map[string]interface{}{
        "event": "job_status_updated",
        "booking_id": bookingID,
        "status": newStatus,
    })
    // if job completed, trigger payment & review flow
    if newStatus == string(domain.StatusCompleted) {
        u.completeJob(ctx, bookingID)
    }
    return nil
}

func (u *JobUsecase) insertStatusTimeline(bookingID int64, status string) {
    // implement DB insert into job_status_updates
}

func (u *JobUsecase) completeJob(ctx context.Context, bookingID int64) {
    // update booking completed_at, calculate final price, call stripe payout
}

func isValidTransition(current, new string) bool {
    transitions := map[string][]string{
        string(domain.StatusAssigned):       {string(domain.StatusHeadingToPickup), string(domain.StatusCancelled)},
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