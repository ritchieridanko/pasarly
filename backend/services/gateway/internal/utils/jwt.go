package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/models"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
)

func JWTParse(token, secret string) (*models.Claim, error) {
	t, err := jwt.ParseWithClaims(
		token,
		&models.Claim{},
		func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claim, ok := t.Claims.(*models.Claim)
	if !ok {
		return nil, ce.ErrInvalidToken
	}

	return claim, nil
}
