package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/ritchieridanko/pasarly/auth-service/configs"
	"github.com/ritchieridanko/pasarly/auth-service/internal/app/models"
	"github.com/ritchieridanko/pasarly/auth-service/internal/app/repositories"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/database"
	"github.com/ritchieridanko/pasarly/auth-service/internal/service/jwt"
	"github.com/ritchieridanko/pasarly/auth-service/internal/shared/ce"
	"github.com/ritchieridanko/pasarly/auth-service/internal/shared/utils"
	"go.opentelemetry.io/otel"
)

const sessionErrTracer string = "usecase.session"

type SessionUsecase interface {
	CreateSession(ctx context.Context, auth *models.Auth, rm *models.RequestMeta) (at *models.AuthToken, err *ce.Error)
}

type sessionUsecase struct {
	config     *configs.Auth
	sr         repositories.SessionRepository
	transactor *database.Transactor
	jwt        *jwt.JWT
}

func NewSessionUsecase(
	cfg *configs.Auth,
	sr repositories.SessionRepository,
	tx *database.Transactor,
	j *jwt.JWT,
) SessionUsecase {
	return &sessionUsecase{config: cfg, sr: sr, transactor: tx, jwt: j}
}

func (u *sessionUsecase) CreateSession(ctx context.Context, auth *models.Auth, rm *models.RequestMeta) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(sessionErrTracer).Start(ctx, "CreateSession")
	defer span.End()

	now := time.Now().UTC()
	sessionToken := utils.NewUUID().String()

	accessToken, errJWT := u.jwt.Create(auth.ID, auth.Role, auth.IsVerified, &now)
	if errJWT != nil {
		e := fmt.Errorf("failed to create session: %w", errJWT)
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
