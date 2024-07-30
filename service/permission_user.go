package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPermissionService struct {
}

func NewUserPermissionService() *UserPermissionService {
	return &UserPermissionService{}
}

// CreateUserPermission 创建新的用户权限关联
func (s *UserPermissionService) CreateUserPermission(ctx *app.Context, userPermission *model.ReqPermissionUserCreate) error {

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now().Format("2006-01-02 15:04:05")
		// 先删除已有的用户权限关联
		err := tx.Where("user_uuid = ?", userPermission.UserUuid).Delete(&model.UserPermission{}).Error
		if err != nil {
			ctx.Logger.Error("Failed to delete user permission by user UUID", err)
			return errors.New("failed to delete user permission by user UUID")
		}

		rlist := make([]*model.UserPermission, 0)
		for _, permissionUuid := range userPermission.PermissionUuids {
			userPermission := &model.UserPermission{
				Uuid:           uuid.New().String(),
				UserUuid:       userPermission.UserUuid,
				PermissionUuid: permissionUuid,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			rlist = append(rlist, userPermission)
		}
		err = tx.Create(&rlist).Error
		if err != nil {
			ctx.Logger.Error("Failed to create user permission", err)
			return errors.New("failed to create user permission")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetUserPermissionByUUID 根据UUID获取用户权限关联
func (s *UserPermissionService) GetUserPermissionByUUID(ctx *app.Context, uuid string) (*model.UserPermission, error) {
	userPermission := &model.UserPermission{}
	err := ctx.DB.Where("uuid = ?", uuid).First(userPermission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user permission not found")
		}
		ctx.Logger.Error("Failed to get user permission by UUID", err)
		return nil, errors.New("failed to get user permission by UUID")
	}
	return userPermission, nil
}

// UpdateUserPermission 更新用户权限关联信息
func (s *UserPermissionService) UpdateUserPermission(ctx *app.Context, userPermission *model.UserPermission) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	userPermission.UpdatedAt = now
	err := ctx.DB.Where("uuid = ?", userPermission.Uuid).Updates(userPermission).Error
	if err != nil {
		ctx.Logger.Error("Failed to update user permission", err)
		return errors.New("failed to update user permission")
	}

	return nil
}

// DeleteUserPermission 删除用户权限关联
func (s *UserPermissionService) DeleteUserPermission(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.UserPermission{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete user permission", err)
		return errors.New("failed to delete user permission")
	}

	return nil
}

// GetUserPermissionList 根据查询参数获取用户权限关联列表
func (s *UserPermissionService) GetUserPermissionList(ctx *app.Context, params *model.ReqUserPermissionQueryParam) (*model.PagedResponse, error) {
	var (
		userPermissions []*model.UserPermission
		total           int64
	)

	db := ctx.DB.Model(&model.UserPermission{})

	if params.UserUuid != "" {
		db = db.Where("user_uuid = ?", params.UserUuid)
	}

	if params.PermissionUuid != "" {
		db = db.Where("permission_uuid = ?", params.PermissionUuid)
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get user permission count", err)
		return nil, errors.New("failed to get user permission count")
	}

	err = db.Offset(params.GetOffset()).Limit(params.PageSize).Find(&userPermissions).Error
	if err != nil {
		ctx.Logger.Error("Failed to get user permission list", err)
		return nil, errors.New("failed to get user permission list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  userPermissions,
	}, nil
}

// 根据用户uuid 获取用户权限关联
func (s *UserPermissionService) GetUserPermissionByUserUUID(ctx *app.Context, userUuid string) ([]*model.UserPermission, error) {
	var userPermissions []*model.UserPermission
	err := ctx.DB.Where("user_uuid = ?", userUuid).Find(&userPermissions).Error
	if err != nil {
		ctx.Logger.Error("Failed to get user permission by user UUID", err)
		return nil, errors.New("failed to get user permission by user UUID")
	}
	return userPermissions, nil
}
