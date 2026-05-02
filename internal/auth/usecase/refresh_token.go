package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type RefreshTokenUsecase struct {
	repo RefreshTokenRepository
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, user_id int64, hash_token string, expires_at time.Time) error
	FindByHashToken(ctx context.Context, hash_token string) (string, error)
	Delete(ctx context.Context, hash_token string) error
}

func NewRefreshTokenUsecase(repo RefreshTokenRepository) *RefreshTokenUsecase {
	return &RefreshTokenUsecase{repo: repo}
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (u *RefreshTokenUsecase) Create(ctx context.Context, user_id int64) (string, error) {
	token, err := generateRefreshToken()
	if err != nil {
		return "", err
	}
	// hash := Hash(token)
	hash_bytes := sha256.Sum256([]byte(token))
	hash := hex.EncodeToString(hash_bytes[:])
	expires_at := time.Now().Add(7 * 24 * time.Hour)
	err = u.repo.Save(ctx, user_id, hash, expires_at)
	return token, err
}
