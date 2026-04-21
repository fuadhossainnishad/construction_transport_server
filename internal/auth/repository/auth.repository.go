package repository

import (
	"construction_transport_server/internal/auth/domain"
)

type AuthRepository interface {
	CreateAccount(user *domain.User) error
	GetAccount(email string) (*domain.User, error)
	VerifyEmail(email string) error
	VerifyOtp(email string, otp string) error
	VerifyPasswordResetToken(email string, token string) error
	UpdatePassword(token string, newPassword string) error
	GetRefreshToken(userID int64) (string, error)
	SetRefreshToken(userID int64, refreshToken string) error
	DeleteRefreshToken(userID int64) error
}
