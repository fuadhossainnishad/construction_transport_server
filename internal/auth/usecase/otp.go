package usecase

import (
	"construction_transport_server/infrastructure/cache/redis"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type OTPService struct {
	store *redis.OTPStore
	pub   EventPublisher
}

type EventPublisher interface {
	PublishUserRegistered(ctx context.Context, email, otp string)
}

func NewOTPService(store *redis.OTPStore, pub EventPublisher) *OTPService {
	return &OTPService{store: store, pub: pub}
}

func (s *OTPService) GenerateAndSendOTP(ctx context.Context, email string) error {
	otp := generateOTP()
	if err := s.store.SetOTP(ctx, email, otp); err != nil {
		return err
	}
	s.pub.PublishUserRegistered(ctx, email, otp)
	return nil
}

func (s *OTPService) VerifyOTP(ctx context.Context, email, otp string) error {
	stored, err := s.store.GetOTP(ctx, email)
	if err != nil {
		return fmt.Errorf("otp not found or expired")
	}
	if stored != otp {
		return fmt.Errorf("invalid otp")
	}
	return s.store.DeleteOTP(ctx, email)
}

func generateOTP() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)
}
