package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type UserRoleService struct {
}

func NewUserRoleService() *UserRoleService {
	return &UserRoleService{}
}

func (s *UserRoleService) CreateUserRole(ctx *app.Context, userRole *model.UserRole) error {
	userRole.CreatedAt = time.Now()
	userRole.UpdatedAt = userRole.CreatedAt

	err := ctx.DB.Create(userRole).Error
	if err != nil {
		ctx.Logger.Error("Failed to create user role", err)
		return errors.New("failed to create user role")
	}
	return nil
}

func (s *UserRoleService) GetUserRoleByUUID(ctx *app.Context, uuid string) (*model.UserRole, error) {
	userRole := &model.UserRole{}
	err := ctx.DB.Where("uuid = ?", uuid).First(userRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user role not found")
		}
		ctx.Logger.Error("Failed to get user role by UUID", err)
		return nil, errors.New("failed to get user role by UUID")
	}
	return userRole, nil
}

func (s *UserRoleService) UpdateUserRole(ctx *app.Context, userRole *model.UserRole) error {
	userRole.UpdatedAt = time.Now()
	err := ctx.DB.Save(userRole).Error
	if err != nil {
		ctx.Logger.Error("Failed to update user role", err)
		return errors.New("failed to update user role")
	}

	return nil
}

func (s *UserRoleService) DeleteUserRole(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.UserRole{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete user role", err)
		return errors.New("failed to delete user role")
	}

	return nil
}

// 获取用户的角色信息
func (s *UserRoleService) GetUserRoleByUserID(ctx *app.Context, userID string) ([]*model.Role, error) {
	var userRole []*model.UserRole
	err := ctx.DB.Where("user_id = ?", userID).Find(&userRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user role not found")
		}
		ctx.Logger.Error("Failed to get user role by user ID", err)
		return nil, errors.New("failed to get user role by user ID")
	}

	var roleIDs []string
	for _, v := range userRole {
		roleIDs = append(roleIDs, v.RoleUUID)
	}

	var roles []*model.Role
	err = ctx.DB.Where("uuid in ?", roleIDs).Find(&roles).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("role not found")
		}
		ctx.Logger.Error("Failed to get role by UUID", err)
		return nil, errors.New("failed to get role by UUID")
	}

	return roles, nil
}
