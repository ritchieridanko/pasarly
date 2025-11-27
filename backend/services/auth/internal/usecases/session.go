package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ritchieridanko/pasarly/backend/services/auth/configs"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/repositories"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
)

const sessionErrTracer string = "usecase.session"

type SessionUsecase interface {
	CreateSession(ctx context.Context, auth *models.Auth, rm *models.RequestMeta) (at *models.AuthToken, err *ce.Error)
	RevokeSession(ctx context.Context, sessionToken string) (err *ce.Error)
}

type sessionUsecase struct {
	config     *configs.Auth
	sr         repositories.SessionRepository
	transactor *database.Transactor
	jwt        *utils.JWT
	validator  *utils.Validator
}

func NewSessionUsecase(
	cfg *configs.Auth,
	sr repositories.SessionRepository,
	tx *database.Transactor,
	j *utils.JWT,
	v *utils.Validator,
) SessionUsecase {
	return &sessionUsecase{config: cfg, sr: sr, transactor: tx, jwt: j, validator: v}
}

func (u *sessionUsecase) CreateSession(ctx context.Context, auth *models.Auth, rm *models.RequestMeta) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(sessionErrTracer).Start(ctx, "CreateSession")
	defer span.End()

	now := time.Now().UTC()
	sessionToken := utils.NewUUID().String()

	accessToken, errJ := u.jwt.Create(auth.ID, auth.Role, auth.IsVerified, &now)
	if errJ != nil {
		e := fmt.Errorf("failed to create session: %w", errJ)
		return nil, ce.NewError(span, ce.CodeJWTCreationFailed, ce.MsgInternalServer, e)
	}

	data := models.CreateSession{
		Token:     sessionToken,
		UserAgent: rm.UserAgent,
		IPAddress: rm.IPAddress,
		ExpiresAt: now.Add(u.config.Token.Duration.Session),
	}

	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		sessionID, err := u.sr.RevokeActiveSession(ctx, auth.ID, rm)
		if err != nil {
			return err
		}
		if sessionID != 0 {
			data.ParentID = &sessionID
		}

		return u.sr.CreateSession(ctx, auth.ID, &data)
	})

	return &models.AuthToken{Session: sessionToken, Access: accessToken}, err
}

func (u *sessionUsecase) RevokeSession(ctx context.Context, sessionToken string) *ce.Error {
	ctx, span := otel.Tracer(sessionErrTracer).Start(ctx, "RevokeSession")
	defer span.End()

	// Validation
	if ok, why := u.validator.Token(&sessionToken); !ok {
		err := fmt.Errorf("failed to revoke session: %w", errors.New(why))
		return ce.NewError(span, ce.CodeInvalidPayload, why, err)
	}

	return u.sr.RevokeSessionByToken(ctx, sessionToken)
}
