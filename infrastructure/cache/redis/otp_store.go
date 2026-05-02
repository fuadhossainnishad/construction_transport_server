package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPStore struct {
	client *redis.Client
}

func NewOTPStore(client *redis.Client) *OTPStore {
	return &OTPStore{
		client: client,
	}
}

func (s *OTPStore) SetOTP(ctx context.Context, email string, otp string) error {
	return s.client.Set(ctx, "otp:"+email, otp, 5*time.Minute).Err()
}

func (s *OTPStore) GetOTP(ctx context.Context, email string) (string, error) {
	return s.client.Get(ctx, "otp:"+email).Result()
}

func (s *OTPStore) DeleteOTP(ctx context.Context, email string) error {
	return s.client.Del(ctx, "otp:"+email).Err()
}
