package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type RoleMenuPermissionService struct {
}

func NewRoleMenuPermissionService() *RoleMenuPermissionService {
	return &RoleMenuPermissionService{}
}

func (s *RoleMenuPermissionService) CreateRoleMenuPermission(ctx *app.Context, rmp *model.RoleMenuPermission) error {
	rmp.CreatedAt = time.Now()
	rmp.UpdatedAt = rmp.CreatedAt

	err := ctx.DB.Create(rmp).Error
	if err != nil {
		ctx.Logger.Error("Failed to create role menu permission", err)
		return errors.New("failed to create role menu permission")
	}
	return nil
}

func (s *RoleMenuPermissionService) GetRoleMenuPermissionByUUID(ctx *app.Context, uuid string) (*model.RoleMenuPermission, error) {
	rmp := &model.RoleMenuPermission{}
	err := ctx.DB.Where("uuid = ?", uuid).First(rmp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("role menu permission not found")
		}
		ctx.Logger.Error("Failed to get role menu permission by UUID", err)
		return nil, errors.New("failed to get role menu permission by UUID")
	}
	return rmp, nil
}

func (s *RoleMenuPermissionService) UpdateRoleMenuPermission(ctx *app.Context, rmp *model.RoleMenuPermission) error {
	rmp.UpdatedAt = time.Now()
	err := ctx.DB.Save(rmp).Error
	if err != nil {
		ctx.Logger.Error("Failed to update role menu permission", err)
		return errors.New("failed to update role menu permission")
	}

	return nil
}

func (s *RoleMenuPermissionService) DeleteRoleMenuPermission(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.RoleMenuPermission{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete role menu permission", err)
		return errors.New("failed to delete role menu permission")
	}

	return nil
}
