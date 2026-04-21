package repository

import (
	"construction_transport_server/internal/auth/domain"
)

type AccountRepository interface {
	CreateAccount(user *domain.User) error
	GetAccount(email string) (*domain.User, error)
	UpdateAccount(user *domain.User) error
	DeleteAccount(Id string) error
}
