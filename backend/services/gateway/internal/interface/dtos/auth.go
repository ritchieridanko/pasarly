package dtos

import "time"

type Auth struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpResponse struct {
	AccessToken string `json:"access_token"`
	Auth        Auth   `json:"auth"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
	Auth        Auth   `json:"auth"`
}
