package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/models"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
)

const addressErrTracer string = "repository.address"

type AddressRepository interface {
	CreateAddress(ctx context.Context, data *models.CreateAddress) (address *models.Address, err *ce.Error)
	GetAllAddresses(ctx context.Context, authID int64) (addresses []models.Address, err *ce.Error)
	UpdateAddress(ctx context.Context, data *models.UpdateAddress) (address *models.Address, err *ce.Error)
	DeleteAddress(ctx context.Context, data *models.DeleteAddress) (err *ce.Error)
	HasPrimary(ctx context.Context, authID int64) (exists bool, err *ce.Error)
	SetPrimary(ctx context.Context, data *models.SetPrimaryAddress) (address *models.Address, err *ce.Error)
	UnsetPrimary(ctx context.Context, authID int64) (address *models.Address, err *ce.Error)
	SetLastUpdatedPrimary(ctx context.Context, authID int64) (address *models.Address, err *ce.Error)
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

func (r *addressRepository) GetAllAddresses(ctx context.Context, authID int64) ([]models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "GetAllAddresses")
	defer span.End()

	query := `
		SELECT
			address_id, recipient, phone, label, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postcode, latitude, longitude, created_at, updated_at
		FROM addresses
		WHERE auth_id = $1
		ORDER BY is_primary DESC, updated_at DESC
	`

	rows, err := r.database.QueryAll(ctx, query, authID)
	if err != nil {
		e := fmt.Errorf("failed to fetch all addresses: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}
	defer rows.Close()

	addresses := make([]models.Address, 0)
	for rows.Next() {
		var address models.Address

		err := rows.Scan(
			&address.ID, &address.Recipient, &address.Phone, &address.Label, &address.Notes,
			&address.IsPrimary, &address.Country, &address.Subdivision1, &address.Subdivision2,
			&address.Subdivision3, &address.Subdivision4, &address.Street, &address.Postcode,
			&address.Latitude, &address.Longitude, &address.CreatedAt, &address.UpdatedAt,
		)
		if err != nil {
			e := fmt.Errorf("failed to fetch all addresses: %w", err)
			return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
		}

		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		e := fmt.Errorf("failed to fetch all addresses: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	if len(addresses) == 0 {
		return []models.Address{}, nil
	}

	return addresses, nil
}

func (r *addressRepository) UpdateAddress(ctx context.Context, data *models.UpdateAddress) (*models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "UpdateAddress")
	defer span.End()

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if data.Recipient != nil {
		setClauses = append(setClauses, fmt.Sprintf("recipient = $%d", argPos))
		args = append(args, *data.Recipient)
		argPos++
	}
	if data.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, *data.Phone)
		argPos++
	}
	if data.Label != nil {
		setClauses = append(setClauses, fmt.Sprintf("label = $%d", argPos))
		args = append(args, *data.Label)
		argPos++
	}
	if data.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", argPos))
		args = append(args, *data.Notes)
		argPos++
	}
	if data.Country != nil {
		setClauses = append(setClauses, fmt.Sprintf("country = $%d", argPos))
		args = append(args, *data.Country)
		argPos++
	}
	if data.Subdivision1 != nil {
		setClauses = append(setClauses, fmt.Sprintf("subdivision_1 = $%d", argPos))
		args = append(args, *data.Subdivision1)
		argPos++
	}
	if data.Subdivision2 != nil {
		setClauses = append(setClauses, fmt.Sprintf("subdivision_2 = $%d", argPos))
		args = append(args, *data.Subdivision2)
		argPos++
	}
	if data.Subdivision3 != nil {
		setClauses = append(setClauses, fmt.Sprintf("subdivision_3 = $%d", argPos))
		args = append(args, *data.Subdivision3)
		argPos++
	}
	if data.Subdivision4 != nil {
		setClauses = append(setClauses, fmt.Sprintf("subdivision_4 = $%d", argPos))
		args = append(args, *data.Subdivision4)
		argPos++
	}
	if data.Street != nil {
		setClauses = append(setClauses, fmt.Sprintf("street = $%d", argPos))
		args = append(args, *data.Street)
		argPos++
	}
	if data.Postcode != nil {
		setClauses = append(setClauses, fmt.Sprintf("postcode = $%d", argPos))
		args = append(args, *data.Postcode)
		argPos++
	}
	if data.Latitude != nil && data.Longitude != nil {
		setClauses = append(setClauses,
			fmt.Sprintf("latitude = $%d", argPos),
			fmt.Sprintf("longitude = $%d", argPos+1),
			fmt.Sprintf("location = ST_SetSRID(ST_MakePoint($%d, $%d), 4326)", argPos+1, argPos),
		)
		args = append(args, *data.Latitude, *data.Longitude)
		argPos += 2
	}
	if len(setClauses) == 0 {
		err := fmt.Errorf("failed to update address: %w", ce.ErrNoFieldsToUpdate)
		return nil, ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, data.AddressID, data.AuthID)

	query := fmt.Sprintf(
		`
			UPDATE addresses
			SET %s
			WHERE address_id = $%d AND auth_id = $%d
			RETURNING
				address_id, recipient, phone, label, notes, is_primary, country,
				subdivision_1, subdivision_2, subdivision_3, subdivision_4,
				street, postcode, latitude, longitude, created_at, updated_at
		`,
		strings.Join(setClauses, ", "), argPos, argPos+1,
	)

	row := r.database.QueryRow(ctx, query, args...)

	var address models.Address
	err := row.Scan(
		&address.ID, &address.Recipient, &address.Phone, &address.Label, &address.Notes,
		&address.IsPrimary, &address.Country, &address.Subdivision1, &address.Subdivision2,
		&address.Subdivision3, &address.Subdivision4, &address.Street, &address.Postcode,
		&address.Latitude, &address.Longitude, &address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to update address: %w", err)
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return nil, ce.NewError(span, ce.CodeAddressNotFound, ce.MsgAddressNotFound, e)
		}

		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return &address, nil
}

func (r *addressRepository) DeleteAddress(ctx context.Context, data *models.DeleteAddress) *ce.Error {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "DeleteAddress")
	defer span.End()

	query := "DELETE FROM addresses WHERE address_id = $1 AND auth_id = $2"

	if err := r.database.Execute(ctx, query, data.AddressID, data.AuthID); err != nil {
		e := fmt.Errorf("failed to delete address: %w", err)
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeAddressNotFound, ce.MsgAddressNotFound, e)
		}

		return ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return nil
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

func (r *addressRepository) SetPrimary(ctx context.Context, data *models.SetPrimaryAddress) (*models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "SetPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = TRUE, updated_at = NOW()
		WHERE address_id = $1 AND auth_id = $2
		RETURNING
			address_id, recipient, phone, label, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postcode, latitude, longitude, created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, data.AddressID, data.AuthID)

	var address models.Address
	err := row.Scan(
		&address.ID, &address.Recipient, &address.Phone, &address.Label, &address.Notes,
		&address.IsPrimary, &address.Country, &address.Subdivision1, &address.Subdivision2,
		&address.Subdivision3, &address.Subdivision4, &address.Street, &address.Postcode,
		&address.Latitude, &address.Longitude, &address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to set address primary: %w", err)
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return nil, ce.NewError(span, ce.CodeAddressNotFound, ce.MsgAddressNotFound, e)
		}

		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return &address, nil
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

func (r *addressRepository) SetLastUpdatedPrimary(ctx context.Context, authID int64) (*models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "SetLastUpdatedPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = TRUE, updated_at = NOW()
		WHERE address_id = (
			SELECT address_id FROM addresses WHERE auth_id = $1
			ORDER BY updated_at DESC LIMIT 1
		)
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
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return nil, nil
		}

		e := fmt.Errorf("failed to set last updated primary: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return &address, nil
}
