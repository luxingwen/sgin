package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionService struct {
}

func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

// CreatePermission 创建新的权限
func (s *PermissionService) CreatePermission(ctx *app.Context, permission *model.Permission) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	permission.CreatedAt = now
	permission.UpdatedAt = now
	permission.Uuid = uuid.New().String()

	// 先查询是否存在相同的权限位
	var isExistPermission model.Permission
	err := ctx.DB.Where("name = ? AND parent_uuid = ?", permission.Name, permission.ParentUuid).First(&isExistPermission).Error
	if err == nil && isExistPermission.Id > 0 {
		return errors.New("permission already exists")
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Logger.Error("Failed to get permission by name and parent_uuid", err)
		return errors.New("failed to get permission by name and parent_uuid")
	}

	err = ctx.DB.Create(permission).Error
	if err != nil {
		ctx.Logger.Error("Failed to create permission", err)
		return errors.New("failed to create permission")
	}
	return nil
}

// GetPermissionByUUID 根据UUID获取权限
func (s *PermissionService) GetPermissionByUUID(ctx *app.Context, uuid string) (*model.Permission, error) {
	permission := &model.Permission{}
	err := ctx.DB.Where("uuid = ?", uuid).First(permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("permission not found")
		}
		ctx.Logger.Error("Failed to get permission by UUID", err)
		return nil, errors.New("failed to get permission by UUID")
	}
	return permission, nil
}

// UpdatePermission 更新权限信息
func (s *PermissionService) UpdatePermission(ctx *app.Context, permission *model.Permission) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	permission.UpdatedAt = now
	err := ctx.DB.Where("uuid = ?", permission.Uuid).Updates(permission).Error
	if err != nil {
		ctx.Logger.Error("Failed to update permission", err)
		return errors.New("failed to update permission")
	}

	return nil
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.Permission{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete permission", err)
		return errors.New("failed to delete permission")
	}

	return nil
}

// GetPermissionList 根据查询参数获取权限列表
func (s *PermissionService) GetPermissionList(ctx *app.Context, params *model.ReqPermissionQueryParam) (*model.PagedResponse, error) {
	var (
		permissions []*model.Permission
		total       int64
	)

	db := ctx.DB.Model(&model.Permission{})

	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get permission count", err)
		return nil, errors.New("failed to get permission count")
	}

	err = db.Find(&permissions).Error
	if err != nil {
		ctx.Logger.Error("Failed to get permission list", err)
		return nil, errors.New("failed to get permission list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  permissions,
	}, nil
}
