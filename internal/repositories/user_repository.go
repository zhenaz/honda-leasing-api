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
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	q *query.Query
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		q: query.Use(db),
	}
}

// FindAll implements UserRepository.
func (r *userRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	users, err := r.q.User.WithContext(ctx).
		Where(r.q.User.IsActive.Is(true)).
		Find()

	if err != nil {
		return nil, fmt.Errorf("userRepository FindAll failed: %w", err)
	}
	return users, nil
}


// FindByID implements UserRepository.
func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := r.q.User.WithContext(ctx).
		Where(
			r.q.User.UserID.Eq(id),
			r.q.User.IsActive.Is(true),
		).
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

// Update implements UserRepository.
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.q.User.WithContext(ctx).
		Where(r.q.User.UserID.Eq(user.UserID)).
		Updates(user)
	if err != nil {
		return fmt.Errorf("userRepository Update failed: %w", err)
	}
	return nil
}

// Delete implements UserRepository. SOFT DELETE dengan mengubah is_active menjadi false
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.q.User.WithContext(ctx).
		Where(r.q.User.UserID.Eq(id)).
		Update(r.q.User.IsActive, false)

	if err != nil {
		return fmt.Errorf("userRepository SoftDelete failed: %w", err)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
