package domain

import "time"

type UserState string

const (
	UserStateActive    UserState = "active"
	UserStatePending   UserState = "pending"
	UserStateSuspended UserState = "suspended"
)

type AuthUser struct {
	ID              int64
	Email           string
	PasswordHash    string
	Role            Role
	State           UserState
	CreatedAt       time.Time
	UpdatedAt       time.Time
	EmailVerifiedAt *time.Time
	DeletedAt       *time.Time
	IsDeleted       bool
	IsActive        bool
	IsVerified      bool
}
