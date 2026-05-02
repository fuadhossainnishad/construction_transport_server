package repository

import "construction_transport_server/internal/account/domain"

type AccountRepository interface {
	CreateAccount(user *domain.UserProfile) error
	GetAccount(email string) (*domain.UserProfile, error)
}
