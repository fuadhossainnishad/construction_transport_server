package usecase

import (
	"construction_transport_server/internal/auth/domain"
	"construction_transport_server/internal/auth/repository"
	"context"
	"errors"
)

type RegisterInput struct {
	Email    string
	Password string
	Role     string
}

type PasswordHashFunc interface {
	Hash(password string) (string, error)
	Compare(hash, plain string) bool
}

type OTPServiceInterface interface {
	GenerateAndSendOTP(ctx context.Context, email string) error
}

type RegisteredUsecase struct {
	repo       repository.AuthRepository
	hashFunc   PasswordHashFunc
	otpService OTPServiceInterface
}

func NewRegisteredUsecase(
	repo repository.AuthRepository,
	hashFunc PasswordHashFunc,
	otpService OTPServiceInterface,
) *RegisteredUsecase {
	return &RegisteredUsecase{
		repo:       repo,
		hashFunc:   hashFunc,
		otpService: otpService,
	}
}

func (u *RegisteredUsecase) Execute(ctx context.Context, input RegisterInput) error {
	hashed, err := u.hashFunc.Hash(input.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := &domain.AuthUser{
		Email:        input.Email,
		PasswordHash: hashed,
		Role:         domain.Role(input.Role),
		IsVerified:   false,
		State:        domain.UserStatePending,
		IsActive:     true,
	}
	if err := u.repo.CreateAuth(ctx, user); err != nil {
		return errors.New("failed to create account")
	}

	// Send OTP asynchronously via event
	if err := u.otpService.GenerateAndSendOTP(ctx, input.Email); err != nil {
		// log but don't fail registration
	}
	return nil
}
