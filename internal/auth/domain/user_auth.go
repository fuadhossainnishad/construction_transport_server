package domain

import "time"

type AuthUser struct {
	ID              int64
	Email           string
	PasswordHash    string
	Role            Role
	CreatedAt       time.Time
	UpdatedAt       time.Time
	EmailVerifiedAt *time.Time
	DeletedAt       *time.Time
	IsActive        bool
	IsVerified      bool
}
