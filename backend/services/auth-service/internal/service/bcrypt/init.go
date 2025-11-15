package bcrypt

import (
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/configs"
	"golang.org/x/crypto/bcrypt"
)

type BCrypt struct {
	config *configs.Auth
}

func Init(cfg *configs.Auth) *BCrypt {
	return &BCrypt{config: cfg}
}

func (s *BCrypt) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.config.BCrypt.Cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (s *BCrypt) Validate(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
