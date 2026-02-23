package service

import (
	"context"

	"honda-leasing-api/internal/domain/model"
	repository "honda-leasing-api/internal/repositories"
)

type UserRoleService interface {
	AssignRole(ctx context.Context, userID, roleID, assignedBy int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]*model.Role, error)
}

type userRoleService struct {
	repo repository.UserRoleRepository
}

func NewUserRoleService(repo repository.UserRoleRepository) UserRoleService {
	return &userRoleService{repo: repo}
}

func (s *userRoleService) AssignRole(
	ctx context.Context,
	userID, roleID, assignedBy int64,
) error {
	return s.repo.AssignRole(ctx, userID, roleID, assignedBy)
}

func (s *userRoleService) GetUserRoles(
	ctx context.Context,
	userID int64,
) ([]*model.Role, error) {
	return s.repo.GetUserRoles(ctx, userID)
}
