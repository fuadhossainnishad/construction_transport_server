package repository

import (
    "construction_transport_server/internal/vehicle/domain"
    "context"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type VehicleRepository interface {
    Create(ctx context.Context, v *domain.Vehicle) error
    GetByID(ctx context.Context, id int64) (*domain.Vehicle, error)
    ListByTransporter(ctx context.Context, transporterID int64) ([]domain.Vehicle, error)
    Update(ctx context.Context, id int64, updates map[string]interface{}) error
    Delete(ctx context.Context, id int64) error
}

type vehicleRepo struct {
    db *pgxpool.Pool
}

func NewVehicleRepository(db *pgxpool.Pool) VehicleRepository {
    return &vehicleRepo{db: db}
}

func (r *vehicleRepo) Create(ctx context.Context, v *domain.Vehicle) error {
    query := `INSERT INTO vehicles (transporter_id, truck_type, license_plate, payload_capacity, axle_count, hourly_rate, is_active)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
    return r.db.QueryRow(ctx, query, v.TransporterID, v.TruckType, v.LicensePlate, v.PayloadCap, v.AxleCount, v.HourlyRate, v.IsActive).Scan(&v.ID)
}

func (r *vehicleRepo) GetByID(ctx context.Context, id int64) (*domain.Vehicle, error) {
    query := `SELECT id, transporter_id, truck_type, license_plate, payload_capacity, axle_count, hourly_rate, is_active, created_at, updated_at
              FROM vehicles WHERE id = $1`
    var v domain.Vehicle
    err := r.db.QueryRow(ctx, query, id).Scan(
        &v.ID, &v.TransporterID, &v.TruckType, &v.LicensePlate, &v.PayloadCap,
        &v.AxleCount, &v.HourlyRate, &v.IsActive, &v.CreatedAt, &v.UpdatedAt)
    if err == pgx.ErrNoRows {
        return nil, nil
    }
    return &v, err
}

func (r *vehicleRepo) ListByTransporter(ctx context.Context, transporterID int64) ([]domain.Vehicle, error) {
    query := `SELECT id, transporter_id, truck_type, license_plate, payload_capacity, axle_count, hourly_rate, is_active, created_at, updated_at
              FROM vehicles WHERE transporter_id = $1`
    rows, err := r.db.Query(ctx, query, transporterID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var vehicles []domain.Vehicle
    for rows.Next() {
        var v domain.Vehicle
        err := rows.Scan(&v.ID, &v.TransporterID, &v.TruckType, &v.LicensePlate, &v.PayloadCap,
            &v.AxleCount, &v.HourlyRate, &v.IsActive, &v.CreatedAt, &v.UpdatedAt)
        if err != nil {
            return nil, err
        }
        vehicles = append(vehicles, v)
    }
    return vehicles, nil
}

func (r *vehicleRepo) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
    // build dynamic query – simplified for brevity; use sqlx or manual string
    query := `UPDATE vehicles SET updated_at = NOW()`
    args := []interface{}{}
    i := 1
    for col, val := range updates {
        query += ", " + col + " = $" + string(rune('0'+i))
        args = append(args, val)
        i++
    }
    query += " WHERE id = $" + string(rune('0'+i))
    args = append(args, id)
    _, err := r.db.Exec(ctx, query, args...)
    return err
}

func (r *vehicleRepo) Delete(ctx context.Context, id int64) error {
    _, err := r.db.Exec(ctx, `DELETE FROM vehicles WHERE id = $1`, id)
    return err
}