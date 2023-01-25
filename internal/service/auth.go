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
	AuthRepo
	salt string
}

func NewAuthSevice(postgres AuthRepo, salt string) *AuthService {
	return &AuthService{postgres, salt}
}

func (s *AuthService) CreateUser(ctx context.Context, user UserSingUp) error {
	var err error
	user.Password, err = s.generateHash(user.Password)
	if err != nil {
		return fmt.Errorf("generate hash failed: %w", err)
	}

	err = s.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) generateHash(password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}
	return string(hash.Sum([]byte(s.salt))), nil
}
