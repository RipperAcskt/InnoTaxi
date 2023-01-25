package service

import (
	"context"
	"crypto/sha1"
	"fmt"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrUserDoesNotExists = fmt.Errorf("user does not exists")
	ErrIncorrectPassword = fmt.Errorf("incorrect password")
)

type UserSingUp struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type UserSingIn struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type AuthRepo interface {
	CreateUser(ctx context.Context, user UserSingUp) error
	CheckUserByEmail(ctx context.Context, email string) (*UserSingIn, error)
}
type AuthService struct {
	postgres AuthRepo
	salt     string
}

func NewAuthSevice(postgres AuthRepo, salt string) *AuthService {
	return &AuthService{postgres: postgres, salt: salt}
}

func (s *AuthService) SingUp(ctx context.Context, user UserSingUp) error {
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

func (s *AuthService) SingIn(ctx context.Context, user UserSingIn) (*Token, error) {
	userDB, err := s.postgres.CheckUserByEmail(ctx, user.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("check user by email failed %w", err)
	}

	hash := sha1.New()
	_, err = hash.Write([]byte(user.Password))
	if err != nil {
		return nil, fmt.Errorf("write failed: %w", err)
	}

	if userDB.Password != string(hash.Sum([]byte(s.salt))) {
		return nil, ErrIncorrectPassword
	}

	token, err := NewToken()
	if err != nil {
		return nil, fmt.Errorf("new token failed: %w", err)
	}

	return token, nil
}
