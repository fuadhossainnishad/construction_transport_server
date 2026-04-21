package domain

import "time"

type UserProfile struct {
	ID           int64
	UserId       int64
	Name         string
	ProfileImage string
	Location     string
	PhoneNumber  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
