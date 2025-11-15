package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/app/models"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/infra/cache"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/shared/ce"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/shared/constants"
	"go.opentelemetry.io/otel"
)

const authErrTracer string = "repository.auth"

type AuthRepository interface {
	CreateAuth(ctx context.Context, data *models.CreateAuth) (auth *models.Auth, err *ce.Error)
	GetAuthByEmail(ctx context.Context, email string) (auth *models.Auth, err *ce.Error)
	IsEmailRegistered(ctx context.Context, email string) (exists bool, err *ce.Error)
	IsEmailReserved(ctx context.Context, email string) (exists bool, err *ce.Error)
}

type authRepository struct {
	database *database.Database
	cache    *cache.Cache
}

func NewAuthRepository(db *database.Database, c *cache.Cache) AuthRepository {
	return &authRepository{database: db, cache: c}
}

func (r *authRepository) CreateAuth(ctx context.Context, data *models.CreateAuth) (*models.Auth, *ce.Error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "CreateAuth")
	defer span.End()

	query := `
		INSERT INTO auth (email, password, role)
		VALUES ($1, $2, $3)
		RETURNING auth_id, email, role, is_verified, created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, data.Email, data.Password, data.Role)

	var auth models.Auth
	err := row.Scan(
		&auth.ID, &auth.Email, &auth.Role, &auth.IsVerified,
		&auth.CreatedAt, &auth.UpdatedAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to create auth: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return &auth, nil
}

func (r *authRepository) GetAuthByEmail(ctx context.Context, email string) (*models.Auth, *ce.Error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "GetAuthByEmail")
	defer span.End()

	query := `
		SELECT auth_id, email, password, role, is_verified, created_at, updated_at
		FROM auth
		WHERE email = $1 AND deleted_at IS NULL
	`
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, email)

	var auth models.Auth
	err := row.Scan(
		&auth.ID, &auth.Email, &auth.Password, &auth.Role, &auth.IsVerified,
		&auth.CreatedAt, &auth.UpdatedAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to fetch auth by email: %w", err)
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, e)
		}

		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return &auth, nil
}

func (r *authRepository) IsEmailRegistered(ctx context.Context, email string) (bool, *ce.Error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "IsEmailRegistered")
	defer span.End()

	query := "SELECT 1 FROM auth WHERE email = $1 AND deleted_at IS NULL"
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, email)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return false, nil
		}

		e := fmt.Errorf("failed to check if email is registered: %w", err)
		return false, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return true, nil
}

func (r *authRepository) IsEmailReserved(ctx context.Context, email string) (bool, *ce.Error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "IsEmailReserved")
	defer span.End()

	key := fmt.Sprintf("%s:%s", constants.CachePrefixEmailReservation, email)

	exists, err := r.cache.Exists(ctx, key)
	if err != nil {
		e := fmt.Errorf("failed to check if email is reserved: %w", err)
		return false, ce.NewError(span, ce.CodeCacheQueryExec, ce.MsgInternalServer, e)
	}

	return exists, nil
}
