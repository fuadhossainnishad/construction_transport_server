package usecase

import (
    "construction_transport_server/internal/vehicle/domain"
    "construction_transport_server/internal/vehicle/repository"
    "context"
    "errors"
)

type VehicleUsecase struct {
    repo repository.VehicleRepository
}

func NewVehicleUsecase(repo repository.VehicleRepository) *VehicleUsecase {
    return &VehicleUsecase{repo: repo}
}

func (u *VehicleUsecase) Create(ctx context.Context, transporterID int64, input domain.CreateVehicleInput) (*domain.Vehicle, error) {
    v := &domain.Vehicle{
        TransporterID: transporterID,
        TruckType:     input.TruckType,
        LicensePlate:  input.LicensePlate,
        PayloadCap:    input.PayloadCap,
        AxleCount:     input.AxleCount,
        HourlyRate:    input.HourlyRate,
        IsActive:      true,
    }
    err := u.repo.Create(ctx, v)
    return v, err
}

func (u *VehicleUsecase) GetByID(ctx context.Context, id, transporterID int64) (*domain.Vehicle, error) {
    v, err := u.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if v == nil || v.TransporterID != transporterID {
        return nil, errors.New("vehicle not found")
    }
    return v, nil
}

func (u *VehicleUsecase) ListMyVehicles(ctx context.Context, transporterID int64) ([]domain.Vehicle, error) {
    return u.repo.ListByTransporter(ctx, transporterID)
}

func (u *VehicleUsecase) Update(ctx context.Context, id, transporterID int64, input domain.UpdateVehicleInput) error {
    v, err := u.repo.GetByID(ctx, id)
    if err != nil || v == nil || v.TransporterID != transporterID {
        return errors.New("vehicle not found")
    }
    updates := make(map[string]interface{})
    if input.TruckType != nil {
        updates["truck_type"] = *input.TruckType
    }
    if input.LicensePlate != nil {
        updates["license_plate"] = *input.LicensePlate
    }
    if input.PayloadCap != nil {
        updates["payload_capacity"] = *input.PayloadCap
    }
    if input.AxleCount != nil {
        updates["axle_count"] = *input.AxleCount
    }
    if input.HourlyRate != nil {
        updates["hourly_rate"] = *input.HourlyRate
    }
    if input.IsActive != nil {
        updates["is_active"] = *input.IsActive
    }
    if len(updates) == 0 {
        return nil
    }
    return u.repo.Update(ctx, id, updates)
}

func (u *VehicleUsecase) Delete(ctx context.Context, id, transporterID int64) error {
    v, err := u.repo.GetByID(ctx, id)
    if err != nil || v == nil || v.TransporterID != transporterID {
        return errors.New("vehicle not found")
    }
    return u.repo.Delete(ctx, id)
}