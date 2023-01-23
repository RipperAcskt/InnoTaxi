package service

type Service struct {
	*AuthService
}

func New(postgres AuthRepo, salt string) *Service {
	return &Service{
		AuthService: NewAuthSevice(postgres, salt),
	}
}
