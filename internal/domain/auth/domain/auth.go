package domain

import "time"

type RegisterRequest struct{
	ID   int64
	Name string
	Email string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
	Role string
	IsActive bool
	IsVerified bool
	EmailVerifiedAt 
}