package eventhandlers

import (
    "construction_transport_server/internal/events"
    "context"
    "encoding/json"
    "log"
)

type PaymentHandler struct {
    // Stripe client
}

func NewPaymentHandler() *PaymentHandler {
    return &PaymentHandler{}
}

func (h *PaymentHandler) HandleJobCompleted(ctx context.Context, payload []byte) error {
    var evt events.JobCompletedPayload
    if err := json.Unmarshal(payload, &evt); err != nil {
        return err
    }
    log.Printf("💰 Processing payout for transporter %d amount %.2f", evt.TransporterID, evt.TotalPrice)
    // call Stripe payout
    return nil
}