package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/models"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
)

const addressErrTracer string = "repository.address"

type AddressRepository interface {
	CreateAddress(ctx context.Context, data *models.CreateAddress) (address *models.Address, err *ce.Error)
	HasPrimary(ctx context.Context, authID int64) (exists bool, err *ce.Error)
	UnsetPrimary(ctx context.Context, authID int64) (address *models.Address, err *ce.Error)
}

type addressRepository struct {
	database *database.Database
}

func NewAddressRepository(db *database.Database) AddressRepository {
	return &addressRepository{database: db}
}

func (r *addressRepository) CreateAddress(ctx context.Context, data *models.CreateAddress) (*models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "CreateAddress")
	defer span.End()

	query := `
		INSERT INTO addresses (
			auth_id, recipient, phone, label, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postcode, latitude, longitude, location
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15, ST_SetSRID(ST_MakePoint($15, $14), 4326)
		)
		RETURNING
			address_id, recipient, phone, label, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postcode, latitude, longitude, created_at, updated_at
	`

	row := r.database.QueryRow(
		ctx, query,
		data.AuthID, data.Recipient, data.Phone, data.Label, data.Notes, data.IsPrimary,
		data.Country, data.Subdivision1, data.Subdivision2, data.Subdivision3,
		data.Subdivision4, data.Street, data.Postcode, data.Latitude, data.Longitude,
	)

	var address models.Address
	err := row.Scan(
		&address.ID, &address.Recipient, &address.Phone, &address.Label, &address.Notes,
		&address.IsPrimary, &address.Country, &address.Subdivision1, &address.Subdivision2,
		&address.Subdivision3, &address.Subdivision4, &address.Street, &address.Postcode,
		&address.Latitude, &address.Longitude, &address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to create address: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return &address, nil
}

func (r *addressRepository) HasPrimary(ctx context.Context, authID int64) (bool, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "HasPrimary")
	defer span.End()

	query := "SELECT 1 FROM addresses WHERE auth_id = $1 AND is_primary = TRUE"
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return false, nil
		}

		e := fmt.Errorf("failed to check if user has primary address: %w", err)
		return false, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return true, nil
}

func (r *addressRepository) UnsetPrimary(ctx context.Context, authID int64) (*models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "UnsetPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = FALSE, updated_at = NOW()
		WHERE auth_id = $1 AND is_primary = TRUE
		RETURNING
			address_id, recipient, phone, label, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postcode, latitude, longitude, created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, authID)

	var address models.Address
	err := row.Scan(
		&address.ID, &address.Recipient, &address.Phone, &address.Label, &address.Notes,
		&address.IsPrimary, &address.Country, &address.Subdivision1, &address.Subdivision2,
		&address.Subdivision3, &address.Subdivision4, &address.Street, &address.Postcode,
		&address.Latitude, &address.Longitude, &address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to unset primary address: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return &address, nil
}
