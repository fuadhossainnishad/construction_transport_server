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

type EventPublisher interface {
	PublishUserRegistered(ctx context.Context, email string, otp string)
}

type RegisteredUsecase struct {
	repo      repository.AuthRepository
	hash_func PasswordHashFunc
	// event_publisher EventPublisher
}

func NewRegisteredUsecase(r repository.AuthRepository, hash_func PasswordHashFunc) *RegisteredUsecase {
	return &RegisteredUsecase{
		repo:      r,
		hash_func: hash_func,
		// event_publisher: event_publisher,
	}
}

func (register_usecase *RegisteredUsecase) Execute(ctx context.Context, input RegisterInput) error {
	hashed_password, err := register_usecase.hash_func.Hash(input.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := &domain.AuthUser{
		Email:        input.Email,
		PasswordHash: hashed_password,
		Role:         domain.Role(input.Role),
		IsVerified:   false,
	}
	err = register_usecase.repo.CreateAuth(ctx, user)
	if err != nil {
		return errors.New("Failed to create account, try again later")
	}

	return nil
}
