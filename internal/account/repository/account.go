package repository

import (
	"construction_transport_server/internal/account/domain"
)

type AccountRepository interface {
	CreateAccount(profile *domain.UserProfile) error
	GetAccount(userID int64) (*domain.UserProfile, error) // change from email to userID
	UpdateAccount(profile *domain.UserProfile) error
}
