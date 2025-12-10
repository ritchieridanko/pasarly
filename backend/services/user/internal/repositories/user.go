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

const userErrTracer string = "repository.user"

type UserRepository interface {
	CreateUser(ctx context.Context, data *models.CreateUser) (user *models.User, err *ce.Error)
	UpsertUser(ctx context.Context, data *models.UpsertUser) (user *models.User, err *ce.Error)
	UpdateUser(ctx context.Context, data *models.UpdateUser) (user *models.User, err *ce.Error)
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

func (r *userRepository) UpsertUser(ctx context.Context, data *models.UpsertUser) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "UpsertUser")
	defer span.End()

	query := `
		INSERT INTO users
			(auth_id, user_id, name, bio, sex, birthdate, phone)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (auth_id)
		DO UPDATE SET
			name = EXCLUDED.name,
			bio = EXCLUDED.bio,
			sex = EXCLUDED.sex,
			birthdate = EXCLUDED.birthdate,
			phone = EXCLUDED.phone,
			updated_at = NOW()
		RETURNING
			user_id, name, bio, sex, birthdate, phone, profile_picture,
			created_at, updated_at
	`

	row := r.database.QueryRow(
		ctx, query,
		data.AuthID, data.UserID, data.Name, data.Bio, data.Sex, data.Birthdate, data.Phone,
	)

	var user models.User
	err := row.Scan(
		&user.ID, &user.Name, &user.Bio, &user.Sex,
		&user.Birthdate, &user.Phone, &user.ProfilePicture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to upsert user: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExec, ce.MsgInternalServer, e)
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, data *models.UpdateUser) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "UpdateUser")
	defer span.End()

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if data.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *data.Name)
		argPos++
	}
	if data.Bio != nil {
		setClauses = append(setClauses, fmt.Sprintf("bio = $%d", argPos))
		args = append(args, *data.Bio)
		argPos++
	}
	if data.Sex != nil {
		setClauses = append(setClauses, fmt.Sprintf("sex = $%d", argPos))
		args = append(args, *data.Sex)
		argPos++
	}
	if data.Birthdate != nil {
		setClauses = append(setClauses, fmt.Sprintf("birthdate = $%d", argPos))
		args = append(args, *data.Birthdate)
		argPos++
	}
	if data.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, *data.Phone)
		argPos++
	}
	if len(setClauses) == 0 {
		err := fmt.Errorf("failed to update user: %w", ce.ErrNoFieldsToUpdate)
		return nil, ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, data.AuthID)

	query := fmt.Sprintf(
		`
			UPDATE users
			SET %s
			WHERE auth_id = $%d AND deleted_at IS NULL
			RETURNING
				user_id, name, bio, sex, birthdate, phone, profile_picture,
				created_at, updated_at
		`,
		strings.Join(setClauses, ", "), argPos,
	)

	row := r.database.QueryRow(ctx, query, args...)

	var user models.User
	err := row.Scan(
		&user.ID, &user.Name, &user.Bio, &user.Sex,
		&user.Birthdate, &user.Phone, &user.ProfilePicture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		e := fmt.Errorf("failed to update user: %w", err)
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, e)
		}

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
