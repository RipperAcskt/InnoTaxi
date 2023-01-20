package service

import (
	"github.com/RipperAcskt/innotaxi/internal/model"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
)

type SingUp interface {
	CreateUser(user model.UserSingUp) error
}

type Service struct {
	SingUp
}

func New(postgres *postgres.Postgres) *Service {
	return &Service{
		SingUp: NewAuthSevice(postgres),
	}
}
