package models

import "time"

type Auth struct {
	ID                int64
	Email             string
	Password          *string
	Role              string
	IsVerified        bool
	EmailChangedAt    *time.Time
	PasswordChangedAt *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateAuth struct {
	Email    string
	Password *string
	Role     string
}
