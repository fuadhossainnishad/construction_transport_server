package usecase

import (
    "construction_transport_server/internal/account/domain"
    "construction_transport_server/internal/account/repository"
    "context"
)

type AccountUsecase struct {
    repo repository.AccountRepository
}

func NewAccountUsecase(repo repository.AccountRepository) *AccountUsecase {
    return &AccountUsecase{repo: repo}
}

func (u *AccountUsecase) GetProfile(ctx context.Context, userID int64) (*domain.UserProfile, error) {
    return u.repo.GetAccount(userID)
}

func (u *AccountUsecase) CreateOrUpdateProfile(ctx context.Context, userID int64, name, phone, image, location string) (*domain.UserProfile, error) {
    existing, _ := u.repo.GetAccount(userID)
    if existing == nil {
        profile := &domain.UserProfile{
            UserId:       userID,
            Name:         name,
            PhoneNumber:  phone,
            ProfileImage: image,
            Location:     location,
        }
        err := u.repo.CreateAccount(profile)
        return profile, err
    }
    // update existing
    if name != "" {
        existing.Name = name
    }
    if phone != "" {
        existing.PhoneNumber = phone
    }
    if image != "" {
        existing.ProfileImage = image
    }
    if location != "" {
        existing.Location = location
    }
    err := u.repo.UpdateAccount(existing)
    return existing, err
}