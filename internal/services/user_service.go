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
	UpdateUser(ctx context.Context, id int64, u *model.User) error
	DeleteUser(ctx context.Context, id int64) error
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

func (s *userService) UpdateUser(ctx context.Context, id int64, u *model.User) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// update fields yang boleh diubah
	existing.Username = u.Username
	existing.Email = u.Email
	existing.FullName = u.FullName
	existing.PhoneNumber = u.PhoneNumber

	// kalau password dikirim, hash ulang
	if u.Password != nil && *u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword(
			[]byte(*u.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}
		hashed := string(hash)
		existing.Password = &hashed
	}

	return s.repo.Update(ctx, existing)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	// optional: cek dulu ada atau tidak
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
