package usecase

import (
	"construction_transport_server/internal/auth/domain"
	"construction_transport_server/internal/auth/repository"

	"context"
	"errors"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

type JWT interface {
	GenerateAccessToken(user_id int64, role domain.Role, email string) (string, error)
}

type LoginUseCase struct {
	repo            repository.AuthRepository
	jwt             JWT
	refresh_usecase *RefreshTokenUsecase
	hash_func       PasswordHashFunc
}

type PasswordHashFunc interface {
	Hash(password string) (string, error)
	Compare(hash_password, raw_password string) bool
}

func NewLoginUseCase(repo repository.AuthRepository, j JWT, refresh_usecase *RefreshTokenUsecase, hash_func PasswordHashFunc) *LoginUseCase {
	return &LoginUseCase{
		repo:            repo,
		jwt:             j,
		refresh_usecase: refresh_usecase,
		hash_func:       hash_func,
	}
}

func (login_usecase *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginResponse, error) {
	user, err := login_usecase.repo.GetAuth(ctx, input.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !login_usecase.hash_func.Compare(user.PasswordHash, input.Password) {
		return nil, errors.New("invalid password")
	}

	access_token, err := login_usecase.jwt.GenerateAccessToken(user.ID, user.Role, user.Email)
	if err != nil {
		return nil, errors.New("security token generation failed, try again later")
	}

	refresh_token, err := login_usecase.refresh_usecase.Create(ctx, user.ID)
	if err != nil {
		return nil, errors.New("refresh token generation failed, try again later")
	}

	return &LoginResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}
