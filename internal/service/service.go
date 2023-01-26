package service

import "github.com/RipperAcskt/innotaxi/config"

type Service struct {
	*AuthService
}

func New(postgres AuthRepo, redis TokenRepo, salt string, cfg *config.Config) *Service {
	return &Service{
		AuthService: NewAuthSevice(postgres, redis, salt, cfg),
	}
}
