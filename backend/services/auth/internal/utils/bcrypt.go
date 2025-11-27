package utils

import (
	"github.com/ritchieridanko/pasarly/backend/services/auth/configs"
	"golang.org/x/crypto/bcrypt"
)

type BCrypt struct {
	config *configs.Auth
}

func NewBCrypt(cfg *configs.Auth) *BCrypt {
	return &BCrypt{config: cfg}
}

func (u *BCrypt) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), u.config.BCrypt.Cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (u *BCrypt) Validate(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
