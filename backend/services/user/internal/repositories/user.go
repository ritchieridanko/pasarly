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

const userErrTracer string = "repository.user"

type UserRepository interface {
	CreateUser(ctx context.Context, data *models.CreateUser) (user *models.User, err *ce.Error)
	Exists(ctx context.Context, authID int64) (exists bool, err *ce.Error)
}

type userRepository struct {
	database *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{database: db}
}

func (r *userRepository) CreateUser(ctx context.Context, data *models.CreateUser) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "CreateUser")
	defer span.End()

	query := `
		INSERT INTO users (auth_id, user_id, name)
		VALUES ($1, $2, $3)
		RETURNING
			user_id, name, bio, sex, birthdate, phone, profile_picture,
			created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, data.AuthID, data.UserID, data.Name)

	var user models.User
	err := row.Scan(
		&user.ID, &user.Name, &user.Bio, &user.Sex,
		&user.Birthdate, &user.Phone, &user.ProfilePicture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to create user: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return &user, nil
}

func (r *userRepository) Exists(ctx context.Context, authID int64) (bool, *ce.Error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "Exists")
	defer span.End()

	query := "SELECT 1 FROM users WHERE auth_id = $1"
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return false, nil
		}

		e := fmt.Errorf("failed to check if user exists: %w", err)
		return false, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return true, nil
}
