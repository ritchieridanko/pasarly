package service

import (
	"github.com/ritchieridanko/pasarly/auth-service/configs"
	"github.com/ritchieridanko/pasarly/auth-service/internal/service/bcrypt"
	"github.com/ritchieridanko/pasarly/auth-service/internal/service/jwt"
)

type Service struct {
	bcrypt *bcrypt.BCrypt
	jwt    *jwt.JWT
}

func Init(cfg *configs.Config) *Service {
	b := bcrypt.Init(&cfg.Auth)
	j := jwt.Init(&cfg.Auth)
	return &Service{bcrypt: b, jwt: j}
}

func (s *Service) BCrypt() *bcrypt.BCrypt {
	return s.bcrypt
}

func (s *Service) JWT() *jwt.JWT {
	return s.jwt
}
