package repository

import (
	"context"
	"fmt"
	"honda-leasing-api/internal/domain/model"
	"honda-leasing-api/internal/domain/query"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type userRepository struct {
	q *query.Query
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		q: query.Use(db),
	}
}

func (r *userRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	users, err := r.q.User.WithContext(ctx).Find()
	if err != nil {
		return nil, fmt.Errorf("userRepository FindAll failed: %w", err)
	}
	return users, nil
}

// FindByID implements UserRepository.
func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := r.q.User.WithContext(ctx).
		Where(r.q.User.UserID.Eq(id)).
		First()
	if err != nil {
		return nil, fmt.Errorf("userRepository FindByID failed: %w", err)
	}
	return user, nil
}

// Create implements UserRepository.
func (d *userRepository) Create(ctx context.Context, user *model.User) error {
	return d.q.User.WithContext(ctx).
		UnderlyingDB().
		Table("account.users").
		Create(user).Error
}