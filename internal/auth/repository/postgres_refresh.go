package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type refreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Save(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *refreshTokenRepo) FindByHashToken(ctx context.Context, tokenHash string) (int64, error) {
	var userID int64
	query := `SELECT user_id FROM refresh_tokens WHERE token_hash = $1 AND expires_at > NOW()`
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *refreshTokenRepo) Delete(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}
