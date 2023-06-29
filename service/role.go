package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type RoleService struct {
}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (s *RoleService) CreateRole(ctx *app.Context, role *model.Role) error {
	role.CreatedAt = time.Now()
	role.UpdatedAt = role.CreatedAt

	err := ctx.DB.Create(role).Error
	if err != nil {
		ctx.Logger.Error("Failed to create role", err)
		return errors.New("failed to create role")
	}
	return nil
}

func (s *RoleService) GetRoleByUUID(ctx *app.Context, uuid string) (*model.Role, error) {
	role := &model.Role{}
	err := ctx.DB.Where("uuid = ?", uuid).First(role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("role not found")
		}
		ctx.Logger.Error("Failed to get role by UUID", err)
		return nil, errors.New("failed to get role by UUID")
	}
	return role, nil
}

func (s *RoleService) UpdateRole(ctx *app.Context, role *model.Role) error {
	role.UpdatedAt = time.Now()
	err := ctx.DB.Save(role).Error
	if err != nil {
		ctx.Logger.Error("Failed to update role", err)
		return errors.New("failed to update role")
	}

	return nil
}

func (s *RoleService) DeleteRole(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.Role{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete role", err)
		return errors.New("failed to delete role")
	}

	return nil
}
