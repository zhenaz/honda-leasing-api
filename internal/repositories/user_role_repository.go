package repository

import (
	"context"
	"errors"

	"honda-leasing-api/internal/domain/model"
	"honda-leasing-api/internal/domain/query"

	"gorm.io/gorm"
)

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID int64, roleID int64, assignedBy int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]*model.Role, error)
}

type userRoleRepository struct {
	q *query.Query
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{
		q: query.Use(db),
	}
}

func (r *userRoleRepository) AssignRole(
	ctx context.Context,
	userID int64,
	roleID int64,
	assignedBy int64,
) error {

	userRole := &model.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: &assignedBy,
	}

	return r.q.UserRole.WithContext(ctx).Create(userRole)
}

func (r *userRoleRepository) GetUserRoles(
	ctx context.Context,
	userID int64,
) ([]*model.Role, error) {

	roles, err := r.q.Role.WithContext(ctx).
		Join(r.q.UserRole,
			r.q.UserRole.RoleID.EqCol(r.q.Role.RoleID)).
		Where(r.q.UserRole.UserID.Eq(userID)).
		Find()

	if err != nil {
		return nil, errors.New("failed get roles")
	}

	return roles, nil
}
