package models

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	AuthID     int64
	Role       string
	IsVerified bool
	jwt.RegisteredClaims
}
