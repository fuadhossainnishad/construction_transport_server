package repository

import (
	"construction_transport_server/internal/auth/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	CreateAuth(ctx context.Context, user *domain.AuthUser) error
	GetAuth(ctx context.Context, email string) (*domain.AuthUser, error)
	VerifyEmail(ctx context.Context, email string) error
	VerifyOtp(ctx context.Context, email string, otp string) error
	VerifyPasswordResetToken(ctx context.Context, email string, token string) error
	UpdatePassword(ctx context.Context, token string, newPassword string) error
}

type authRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateAuth(ctx context.Context, user *domain.AuthUser) error {

	query := `
	INSERT INTO auth(
	email,
	password_hash,
	role,
	state,
	crated_at,
	updated_at,
	email_verified_at,
	deleted_at,
	is_deleted,
	is_active,
	is_verified
	)
	VALUES ($1,$2,$3,$4,NOW(),NOW(),NOW)
	`

	return nil
}

func (r *authRepo) GetAuth(ctx context.Context, email string) (*domain.AuthUser, error) {
	return nil, nil
}

func (r *authRepo) VerifyEmail(ctx context.Context, email string) error {
	return nil
}

func (r *authRepo) VerifyOtp(ctx context.Context, email string, otp string) error {
	return nil
}

func (r *authRepo) VerifyPasswordResetToken(ctx context.Context, email string, token string) error {
	return nil
}

func (r *authRepo) UpdatePassword(ctx context.Context, token string, newPassword string) error {
	return nil
}
