package service

import (
	"context"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/model"
)

type Service struct {
	*AuthService
	*UserService
}
type Repo interface {
	AuthRepo
	UserRepo
}
type UserRepo interface {
	GetUserById(ctx context.Context, id string) (*model.User, error)
	UpdateUserById(ctx context.Context, id string, user *model.User) error
	DeleteUserById(ctx context.Context, id string) error
}
type UserService struct {
	UserRepo
}

func New(postgres Repo, redis TokenRepo, salt string, cfg *config.Config) *Service {
	return &Service{
		AuthService: NewAuthSevice(postgres, salt, cfg),
		UserService: NewUserService(postgres),
	}
}

func NewUserService(postgres UserRepo) *UserService {
	return &UserService{postgres}
}

func (user *UserService) GetProfile(ctx context.Context, id string) (*model.User, error) {
	return user.GetUserById(ctx, id)
}

func (user *UserService) UpdateProfile(ctx context.Context, id string, userUpdate *model.User) error {
	return user.UpdateUserById(ctx, id, userUpdate)
}

func (user *UserService) DeleteUser(ctx context.Context, id string) error {
	return user.DeleteUserById(ctx, id)
}
