package models

import "time"

type User struct {
	ID             string
	Name           string
	Bio            *string
	Sex            *string
	Birthdate      *time.Time
	Phone          *string
	ProfilePicture *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateUser struct {
	AuthID int64
	UserID string
	Name   string
}

type UpsertUser struct {
	AuthID    int64
	UserID    string
	Name      string
	Bio       *string
	Sex       *string
	Birthdate *time.Time
	Phone     *string
}

type UpdateUser struct {
	AuthID    int64
	Name      *string
	Bio       *string
	Sex       *string
	Birthdate *time.Time
	Phone     *string
}
