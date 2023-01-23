package service

import (
	"context"
	"crypto/sha1"
	"fmt"
)

var (
	ErrUserAlreadyExists error = fmt.Errorf("user already exists")
)

type UserSingUp struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type AuthRepo interface {
	CreateUser(ctx context.Context, user UserSingUp) error
}
type AuthService struct {
	postgres AuthRepo
	salt     string
}

func NewAuthSevice(postgres AuthRepo, salt string) *AuthService {
	return &AuthService{postgres: postgres, salt: salt}
}

func (s *AuthService) CreateUser(ctx context.Context, user UserSingUp) error {
	var err error
	user.Password, err = s.generateHash(user.Password)
	if err != nil {
		return fmt.Errorf("generate hash failed: %w", err)
	}

	err = s.postgres.CreateUser(ctx, user)
	if err == ErrUserAlreadyExists {
		return ErrUserAlreadyExists
	} else {
		return err
	}
}

func (s *AuthService) generateHash(password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}
	return string(hash.Sum([]byte(s.salt))), nil
}
