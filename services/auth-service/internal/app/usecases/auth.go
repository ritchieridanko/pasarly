package usecases

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/pasarly/auth-service/internal/app/models"
	"github.com/ritchieridanko/pasarly/auth-service/internal/app/repositories"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/database"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/auth-service/internal/service/bcrypt"
	"github.com/ritchieridanko/pasarly/auth-service/internal/shared/ce"
	"github.com/ritchieridanko/pasarly/auth-service/internal/shared/constants"
	"github.com/ritchieridanko/pasarly/auth-service/internal/shared/utils"
	"go.opentelemetry.io/otel"
)

const authErrTracer string = "usecase.auth"

type AuthUsecase interface {
	SignUp(ctx context.Context, data *models.CreateAuth) (auth *models.Auth, err *ce.Error)
}

type authUsecase struct {
	ar         repositories.AuthRepository
	tr         repositories.TokenRepository
	transactor *database.Transactor
	bcrypt     *bcrypt.BCrypt
	logger     *logger.Logger
}

func NewAuthUsecase(
	ar repositories.AuthRepository,
	tr repositories.TokenRepository,
	tx *database.Transactor,
	b *bcrypt.BCrypt,
	l *logger.Logger,
) AuthUsecase {
	return &authUsecase{ar: ar, tr: tr, transactor: tx, bcrypt: b, logger: l}
}

func (u *authUsecase) SignUp(ctx context.Context, data *models.CreateAuth) (*models.Auth, *ce.Error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "SignUp")
	defer span.End()

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

		h, errH := u.bcrypt.Hash(*data.Password)
		if errH != nil {
			e := fmt.Errorf("failed to sign up: %w", errH)
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

	verificationToken := utils.NewUUID().String()
	if err := u.tr.CreateVerificationToken(ctx, auth.ID, verificationToken); err != nil {
		u.logger.Sugar().Warnln(err.Error())
	}

	return auth, nil
}
