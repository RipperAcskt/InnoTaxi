package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/RipperAcskt/innotaxi/internal/model"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
)

var (
	salt                       = "124jkhsdaf3425"
	ErrUserAlreadyExists error = fmt.Errorf("user already exists")
)

type AuthService struct {
	postgres *postgres.Postgres
}

func NewAuthSevcie(postgres *postgres.Postgres) *AuthService {
	return &AuthService{postgres: postgres}
}

func (s *AuthService) CreateUser(user model.UserSingUp) error {
	var err error
	user.Password, err = generateHash(user.Password)
	fmt.Println(user.Password)
	if err != nil {
		return fmt.Errorf("generate hash failed: %v", err)
	}

	err = s.postgres.CreateUser(user)
	if err == postgres.ErrUserAlreadyExists {
		return ErrUserAlreadyExists
	} else {
		return err
	}
}

func generateHash(password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", fmt.Errorf("write failed: %v", err)
	}
	return string(hash.Sum([]byte(salt))), nil
}
