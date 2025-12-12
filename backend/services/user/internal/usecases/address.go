package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/repositories"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
)

const addressErrTracer string = "usecase.address"

type AddressUsecase interface {
	CreateAddress(ctx context.Context, data *models.CreateAddress) (address *models.Address, oldPrimaryAddress *models.Address, err *ce.Error)
	GetAllAddresses(ctx context.Context, authID int64) (addresses []models.Address, err *ce.Error)
	UpdateAddress(ctx context.Context, data *models.UpdateAddress) (address *models.Address, err *ce.Error)
	SetPrimaryAddress(ctx context.Context, data *models.SetPrimaryAddress) (newPrimaryAddress *models.Address, oldPrimaryAddress *models.Address, err *ce.Error)
}

type addressUsecase struct {
	ar         repositories.AddressRepository
	transactor *database.Transactor
	validator  *utils.Validator
}

func NewAddressUsecase(ar repositories.AddressRepository, tx *database.Transactor, v *utils.Validator) AddressUsecase {
	return &addressUsecase{ar: ar, transactor: tx, validator: v}
}

func (u *addressUsecase) CreateAddress(ctx context.Context, data *models.CreateAddress) (*models.Address, *models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "CreateAddress")
	defer span.End()

	// Validations
	if ok, why := u.validator.Name(&data.Recipient, false); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrPhone(&data.Phone, false); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrLabel(&data.Label, false); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrNotes(data.Notes); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrCountry(&data.Country, false); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrSubdivision(data.Subdivision1); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrSubdivision(data.Subdivision2); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrSubdivision(data.Subdivision3); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrSubdivision(data.Subdivision4); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrStreet(&data.Street, false); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrPostcode(&data.Postcode, false); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrLatitude(&data.Latitude, false); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrLongitude(&data.Longitude, false); !ok {
		err := fmt.Errorf("failed to create address: %w", errors.New(why))
		return nil, nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}

	var address, oldPrimaryAddress *models.Address
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		exists, err := u.ar.HasPrimary(ctx, data.AuthID)
		if err != nil {
			return err
		}
		if exists && data.IsPrimary {
			oldPrimaryAddress, err = u.ar.UnsetPrimary(ctx, data.AuthID)
			if err != nil {
				return err
			}
		}
		if !exists {
			data.IsPrimary = true
		}

		address, err = u.ar.CreateAddress(ctx, data)
		return err
	})

	return address, oldPrimaryAddress, err
}

func (u *addressUsecase) GetAllAddresses(ctx context.Context, authID int64) ([]models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "GetAllAddresses")
	defer span.End()

	return u.ar.GetAllAddresses(ctx, authID)
}

func (u *addressUsecase) UpdateAddress(ctx context.Context, data *models.UpdateAddress) (*models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "UpdateAddress")
	defer span.End()

	// Validations
	if ok, why := u.validator.Name(data.Recipient, true); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrPhone(data.Phone, true); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrLabel(data.Label, true); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrNotes(data.Notes); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrCountry(data.Country, true); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrSubdivision(data.Subdivision1); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrSubdivision(data.Subdivision2); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrSubdivision(data.Subdivision3); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrSubdivision(data.Subdivision4); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrStreet(data.Street, true); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrPostcode(data.Postcode, true); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrLatitude(data.Latitude, true); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.AddrLongitude(data.Longitude, true); !ok {
		err := fmt.Errorf("failed to update address: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}

	return u.ar.UpdateAddress(ctx, data)
}

func (u *addressUsecase) SetPrimaryAddress(ctx context.Context, data *models.SetPrimaryAddress) (*models.Address, *models.Address, *ce.Error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "SetPrimaryAddress")
	defer span.End()

	var newPrimaryAddress, oldPrimaryAddress *models.Address
	err := u.transactor.WithTx(ctx, func(ctx context.Context) (err *ce.Error) {
		oldPrimaryAddress, err = u.ar.UnsetPrimary(ctx, data.AuthID)
		if err != nil {
			return err
		}

		newPrimaryAddress, err = u.ar.SetPrimary(ctx, data)
		return err
	})

	return newPrimaryAddress, oldPrimaryAddress, err
}
