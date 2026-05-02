package repository

type SessionRepository interface {
	GetRefreshToken(userID int64) (string, error)
	SetRefreshToken(userID int64, refreshToken string) error
	DeleteRefreshToken(userID int64) error
}
