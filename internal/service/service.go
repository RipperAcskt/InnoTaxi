package service

import "github.com/RipperAcskt/innotaxi/config"

type Service struct {
	*AuthService
}

func New(postgres AuthRepo, salt string, cfg *config.Config) *Service {
	return &Service{
		AuthService: NewAuthSevice(postgres, salt, cfg),
	}
}
