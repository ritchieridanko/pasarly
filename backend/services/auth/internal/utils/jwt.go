package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/models"
)

type JWT struct {
	issuer   string
	secret   string
	duration time.Duration
}

func NewJWT(issuer, secret string, d time.Duration) *JWT {
	return &JWT{issuer: issuer, secret: secret, duration: d}
}

func (u *JWT) Create(authID int64, role string, isVerified bool, now *time.Time) (string, error) {
	if now == nil {
		t := time.Now().UTC()
		now = &t
	}

	c := models.Claim{
		AuthID:     authID,
		Role:       role,
		IsVerified: isVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    u.issuer,
			Subject:   fmt.Sprintf("%d", authID),
			IssuedAt:  &jwt.NumericDate{Time: *now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(u.duration)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(u.secret))
}
