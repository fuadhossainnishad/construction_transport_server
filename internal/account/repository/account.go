package repository

import (
	"construction_transport_server/internal/account/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type accountRepo struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) CreateAccount(profile *domain.UserProfile) error {
	query := `INSERT INTO profiles (user_id, full_name, phone_number, profile_image, location) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query, profile.UserId, profile.Name, profile.PhoneNumber, profile.ProfileImage, profile.Location)
	return err
}

func (r *accountRepo) GetAccount(userID int64) (*domain.UserProfile, error) {
	query := `SELECT id, user_id, full_name, profile_image, location, phone_number, created_at, updated_at 
              FROM profiles WHERE user_id = $1`
	var profile domain.UserProfile
	err := r.db.QueryRow(context.Background(), query, userID).Scan(
		&profile.ID, &profile.UserId, &profile.Name, &profile.ProfileImage,
		&profile.Location, &profile.PhoneNumber, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *accountRepo) UpdateAccount(profile *domain.UserProfile) error {
	query := `UPDATE profiles SET full_name = $1, phone_number = $2, profile_image = $3, location = $4, updated_at = NOW() WHERE user_id = $5`
	_, err := r.db.Exec(context.Background(), query, profile.Name, profile.PhoneNumber, profile.ProfileImage, profile.Location, profile.UserId)
	return err
}
