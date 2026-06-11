package eventhandlers

import (
	"construction_transport_server/internal/events"
	"context"
	"encoding/json"
	"log"
)

type EmailHandler struct {
	// email service client
}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{}
}

func (h *EmailHandler) HandleUserRegistered(ctx context.Context, payload []byte) error {
	var evt events.UserRegisteredPayload
	if err := json.Unmarshal(payload, &evt); err != nil {
		return err
	}
	log.Printf("📧 Sending OTP %s to %s", evt.OTP, evt.Email)
	// call actual email service here
	return nil
}

func (h *EmailHandler) HandleJobStatusChanged(ctx context.Context, payload []byte) error {
	var evt events.JobStatusChangedPayload
	if err := json.Unmarshal(payload, &evt); err != nil {
		return err
	}
	log.Printf("📧 Sending job status update email for booking %d: %s", evt.BookingID, evt.Status)
	return nil
}
