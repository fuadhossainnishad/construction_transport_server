package repository

import (
	"construction_transport_server/internal/auth/domain"
	"context"
)

type AuthRepository interface {
	CreateAuth(ctx context.Context, user *domain.AuthUser) error
	GetAuth(ctx context.Context, email string) (*domain.AuthUser, error)
	VerifyEmail(ctx context.Context, email string) error
	VerifyOtp(ctx context.Context, email string, otp string) error
	VerifyPasswordResetToken(ctx context.Context, email string, token string) error
	UpdatePassword(ctx context.Context, token string, newPassword string) error
}
