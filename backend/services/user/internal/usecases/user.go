package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/user/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/repositories"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
)

const userErrTracer string = "usecase.user"

type UserUsecase interface {
	UpsertUser(ctx context.Context, data *models.UpsertUser) (user *models.User, err *ce.Error)
	GetUser(ctx context.Context, authID int64) (user *models.User, err *ce.Error)
	UpdateUser(ctx context.Context, data *models.UpdateUser) (user *models.User, err *ce.Error)
}

type userUsecase struct {
	ur        repositories.UserRepository
	validator *utils.Validator
}

func NewUserUsecase(ur repositories.UserRepository, v *utils.Validator) UserUsecase {
	return &userUsecase{ur: ur, validator: v}
}

func (u *userUsecase) UpsertUser(ctx context.Context, data *models.UpsertUser) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "UpsertUser")
	defer span.End()

	// Validations
	if ok, why := u.validator.Name(&data.Name, false); !ok {
		err := fmt.Errorf("failed to upsert user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.Bio(data.Bio); !ok {
		err := fmt.Errorf("failed to upsert user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.Sex(data.Sex); !ok {
		err := fmt.Errorf("failed to upsert user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.Birthdate(data.Birthdate); !ok {
		err := fmt.Errorf("failed to upsert user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.Phone(data.Phone); !ok {
		err := fmt.Errorf("failed to upsert user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}

	data.UserID = utils.NewUUID().String()
	return u.ur.UpsertUser(ctx, data)
}

func (u *userUsecase) GetUser(ctx context.Context, authID int64) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "GetUser")
	defer span.End()

	return u.ur.GetUserByAuthID(ctx, authID)
}

func (u *userUsecase) UpdateUser(ctx context.Context, data *models.UpdateUser) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "UpdateUser")
	defer span.End()

	// Validations
	if ok, why := u.validator.Name(data.Name, true); !ok {
		err := fmt.Errorf("failed to update user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.Bio(data.Bio); !ok {
		err := fmt.Errorf("failed to update user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.Sex(data.Sex); !ok {
		err := fmt.Errorf("failed to update user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.Birthdate(data.Birthdate); !ok {
		err := fmt.Errorf("failed to update user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.Phone(data.Phone); !ok {
		err := fmt.Errorf("failed to update user: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}

	return u.ur.UpdateUser(ctx, data)
}
