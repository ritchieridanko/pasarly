package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ritchieridanko/pasarly/auth-service/internal/app/models"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/database"
	"github.com/ritchieridanko/pasarly/auth-service/internal/shared/ce"
	"go.opentelemetry.io/otel"
)

const sessionErrTracer string = "repository.session"

type SessionRepository interface {
	CreateSession(ctx context.Context, authID int64, data *models.CreateSession) (err *ce.Error)
	RevokeActiveSession(ctx context.Context, authID int64, rm *models.RequestMeta) (sessionID int64, err *ce.Error)
}

type sessionRepository struct {
	database *database.Database
}

func NewSessionRepository(db *database.Database) SessionRepository {
	return &sessionRepository{database: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, authID int64, data *models.CreateSession) *ce.Error {
	ctx, span := otel.Tracer(sessionErrTracer).Start(ctx, "CreateSession")
	defer span.End()

	query := `
		INSERT INTO sessions (auth_id, parent_id, token, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	err := r.database.Execute(
		ctx, query,
		authID, data.ParentID, data.Token, data.UserAgent, data.IPAddress, data.ExpiresAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to create session: %w", err)
		return ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return nil
}

func (r *sessionRepository) RevokeActiveSession(ctx context.Context, authID int64, rm *models.RequestMeta) (int64, *ce.Error) {
	ctx, span := otel.Tracer(sessionErrTracer).Start(ctx, "RevokeActiveSession")
	defer span.End()

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE
			auth_id = $1 AND user_agent = $2 AND ip_address = $3 AND
			revoked_at IS NULL AND expires_at >= $4
		RETURNING session_id
	`

	row := r.database.QueryRow(
		ctx, query,
		authID, rm.UserAgent, rm.IPAddress, time.Now().UTC(),
	)

	var sessionID int64
	if err := row.Scan(&sessionID); err != nil {
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return 0, nil
		}

		e := fmt.Errorf("failed to revoke active session: %w", err)
		return 0, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return sessionID, nil
}
