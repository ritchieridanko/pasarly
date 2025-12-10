package dtos

import "time"

type User struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Bio            *string    `json:"bio"`
	Sex            *string    `json:"sex"`
	Birthdate      *time.Time `json:"birthdate"`
	Phone          *string    `json:"phone"`
	ProfilePicture *string    `json:"profile_picture"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type UpsertUserRequest struct {
	Name      string     `json:"name" binding:"required"`
	Bio       *string    `json:"bio"`
	Sex       *string    `json:"sex"`
	Birthdate *time.Time `json:"birthdate" time_format:"2006-01-02"`
	Phone     *string    `json:"phone"`
}

type UpsertUserResponse struct {
	User User `json:"user"`
}
