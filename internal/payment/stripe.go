package payment

import (
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"
	"github.com/stripe/stripe-go/v79/transfer"
)

type StripeClient struct {
	secretKey string
}

func NewStripeClient(secretKey string) *StripeClient {
	stripe.Key = secretKey
	return &StripeClient{secretKey: secretKey}
}

func (s *StripeClient) CreatePaymentIntent(amount int64, currency string, customerID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Customer: stripe.String(customerID),
	}
	return paymentintent.New(params)
}

func (s *StripeClient) CreateTransfer(amount int64, currency, destination string) (*stripe.Transfer, error) {
	params := &stripe.TransferParams{
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(currency),
		Destination: stripe.String(destination),
	}
	return transfer.New(params)
}
