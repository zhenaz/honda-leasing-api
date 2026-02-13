package service

import (
	"context"
	"errors"
	"honda-leasing-api/internal/domain/model"
	repository "honda-leasing-api/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.FindByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, u *model.User) error {
	if u.Password != nil && *u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword(
			[]byte(*u.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}
		hashed := string(hash)
		u.Password = &hashed
	}

	active := true
	u.IsActive = &active

	return s.repo.Create(ctx, u)
}
