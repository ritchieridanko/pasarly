package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/publisher"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/repositories"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"github.com/ritchieridanko/pasarly/backend/shared/events/v1"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const authErrTracer string = "usecase.auth"

type AuthUsecase interface {
	SignUp(ctx context.Context, data *models.CreateAuth) (auth *models.Auth, err *ce.Error)
	SignIn(ctx context.Context, data *models.GetAuth) (auth *models.Auth, err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (exists bool, err *ce.Error)
}

type authUsecase struct {
	ar         repositories.AuthRepository
	tr         repositories.TokenRepository
	transactor *database.Transactor
	acp        *publisher.Publisher
	bcrypt     *utils.BCrypt
	validator  *utils.Validator
	logger     *logger.Logger
}

func NewAuthUsecase(
	ar repositories.AuthRepository,
	tr repositories.TokenRepository,
	tx *database.Transactor,
	acp *publisher.Publisher,
	b *utils.BCrypt,
	v *utils.Validator,
	l *logger.Logger,
) AuthUsecase {
	return &authUsecase{ar: ar, tr: tr, transactor: tx, acp: acp, bcrypt: b, validator: v, logger: l}
}

func (u *authUsecase) SignUp(ctx context.Context, data *models.CreateAuth) (*models.Auth, *ce.Error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "SignUp")
	defer span.End()

	// Validations
	if ok, why := u.validator.Email(&data.Email); !ok {
		err := fmt.Errorf("failed to sign up: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}
	if ok, why := u.validator.Password(data.Password); !ok {
		err := fmt.Errorf("failed to sign up: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}

	var auth *models.Auth
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		email := utils.NormalizeString(data.Email)
		exists, err := u.ar.IsEmailRegistered(ctx, email)
		if err != nil {
			return err
		}
		if exists {
			e := fmt.Errorf("failed to sign up: %w", ce.ErrEmailAlreadyRegistered)
			return ce.NewError(span, ce.CodeDataConflict, ce.MsgEmailAlreadyRegistered, e)
		}

		exists, err = u.ar.IsEmailReserved(ctx, email)
		if err != nil {
			return err
		}
		if exists {
			e := fmt.Errorf("failed to sign up: %w", ce.ErrEmailReserved)
			return ce.NewError(span, ce.CodeDataConflict, ce.MsgEmailAlreadyRegistered, e)
		}

		h, eh := u.bcrypt.Hash(*data.Password)
		if eh != nil {
			e := fmt.Errorf("failed to sign up: %w", eh)
			return ce.NewError(span, ce.CodeHashingFailed, ce.MsgInternalServer, e)
		}

		ca := models.CreateAuth{
			Email:    email,
			Password: &h,
			Role:     constants.RoleCustomer,
		}

		auth, err = u.ar.CreateAuth(ctx, &ca)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Create and store verification token in cache
	token := utils.NewUUID().String()
	if err := u.tr.CreateVerificationToken(ctx, auth.ID, token); err != nil {
		u.logger.Sugar().Warnln(err.Error())
		return auth, nil
	}

	// Publish event
	key := fmt.Sprintf("auth_%d", auth.ID)
	evt := events.AuthCreated{
		EventId:   utils.NewUUID().String(),
		Email:     auth.Email,
		Token:     token,
		CreatedAt: timestamppb.New(time.Now().UTC()),
	}

	_ = u.acp.Publish(ctx, key, &evt) // failed to publish event does not fail SignUp process

	return auth, nil
}

func (u *authUsecase) SignIn(ctx context.Context, data *models.GetAuth) (*models.Auth, *ce.Error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "SignIn")
	defer span.End()

	// Validation
	if ok, why := u.validator.Email(&data.Email); !ok {
		err := fmt.Errorf("failed to sign in: %w", errors.New(why))
		return nil, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}

	email := utils.NormalizeString(data.Email)
	auth, err := u.ar.GetAuthByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if auth.Password == nil {
		err := fmt.Errorf("failed to sign in: %w", ce.ErrWrongSignInMethod)
		return nil, ce.NewError(span, ce.CodeWrongSignInMethod, ce.MsgInvalidCredentials, err)
	}

	if err := u.bcrypt.Validate(*auth.Password, data.Password); err != nil {
		e := fmt.Errorf("failed to sign in: %w", err)
		return nil, ce.NewError(span, ce.CodeInvalidCredentials, ce.MsgInvalidCredentials, e)
	}

	return auth, nil
}

func (u *authUsecase) IsEmailAvailable(ctx context.Context, email string) (bool, *ce.Error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "IsEmailAvailable")
	defer span.End()

	// Validation
	if ok, why := u.validator.Email(&email); !ok {
		err := fmt.Errorf("failed to check if email is available: %w", errors.New(why))
		return false, ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}

	email = utils.NormalizeString(email)
	exists, err := u.ar.IsEmailRegistered(ctx, email)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	exists, err = u.ar.IsEmailReserved(ctx, email)
	if err != nil {
		return false, err
	}

	return !exists, nil
}
