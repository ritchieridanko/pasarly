package utils

import "golang.org/x/crypto/bcrypt"

type BCrypt struct {
	cost int
}

func NewBCrypt(cost int) *BCrypt {
	return &BCrypt{cost: cost}
}

func (u *BCrypt) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), u.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (u *BCrypt) Validate(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
