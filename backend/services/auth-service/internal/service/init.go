package service

import (
	"github.com/ritchieridanko/pasarly/auth-service/configs"
	"github.com/ritchieridanko/pasarly/auth-service/internal/service/bcrypt"
	"github.com/ritchieridanko/pasarly/auth-service/internal/service/jwt"
	"github.com/ritchieridanko/pasarly/auth-service/internal/service/validator"
)

type Service struct {
	bcrypt    *bcrypt.BCrypt
	jwt       *jwt.JWT
	validator *validator.Validator
}

func Init(cfg *configs.Config) *Service {
	b := bcrypt.Init(&cfg.Auth)
	j := jwt.Init(&cfg.Auth)
	v := validator.Init()
	return &Service{bcrypt: b, jwt: j, validator: v}
}

func (s *Service) BCrypt() *bcrypt.BCrypt {
	return s.bcrypt
}

func (s *Service) JWT() *jwt.JWT {
	return s.jwt
}

func (s *Service) Validator() *validator.Validator {
	return s.validator
}
