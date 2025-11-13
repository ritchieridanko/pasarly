package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/pasarly/auth-service/configs"
	"github.com/ritchieridanko/pasarly/auth-service/internal/app/models"
)

type JWT struct {
	config *configs.Auth
}

func Init(cfg *configs.Auth) *JWT {
	return &JWT{config: cfg}
}

func (s *JWT) Create(authID int64, role string, isVerified bool, now *time.Time) (string, error) {
	if now == nil {
		t := time.Now().UTC()
		now = &t
	}

	c := models.Claim{
		AuthID:     authID,
		Role:       role,
		IsVerified: isVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.JWT.Issuer,
			Subject:   fmt.Sprintf("%d", authID),
			IssuedAt:  &jwt.NumericDate{Time: *now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(s.config.JWT.Duration)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(s.config.JWT.Secret))
}
