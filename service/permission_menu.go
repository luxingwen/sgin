package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionMenuService struct {
}

func NewPermissionMenuService() *PermissionMenuService {
	return &PermissionMenuService{}
}

// CreatePermissionMenu 创建新的权限菜单关联
func (s *PermissionMenuService) CreatePermissionMenu(ctx *app.Context, permissionMenu *model.ReqPermissionMenuCreate) error {

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now().Format("2006-01-02 15:04:05")
		// 先删除已有的权限菜单关联
		err := tx.Where("permission_uuid = ?", permissionMenu.PermissionUuid).Delete(&model.PermissionMenu{}).Error
		if err != nil {
			ctx.Logger.Error("Failed to delete permission menu by menu UUID", err)
			return errors.New("failed to delete permission menu by menu UUID")
		}

		rlist := make([]*model.PermissionMenu, 0)
		for _, menuUuid := range permissionMenu.MenuUuids {
			permissionMenu := &model.PermissionMenu{
				Uuid:           uuid.New().String(),
				PermissionUuid: permissionMenu.PermissionUuid,
				MenuUuid:       menuUuid,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			rlist = append(rlist, permissionMenu)
		}
		err = tx.Create(&rlist).Error
		if err != nil {
			ctx.Logger.Error("Failed to create permission menu", err)
			return errors.New("failed to create permission menu")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// GetPermissionMenuByUUID 根据UUID获取权限菜单关联
func (s *PermissionMenuService) GetPermissionMenuByUUID(ctx *app.Context, uuid string) (*model.PermissionMenu, error) {
	permissionMenu := &model.PermissionMenu{}
	err := ctx.DB.Where("uuid = ?", uuid).First(permissionMenu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("permission menu not found")
		}
		ctx.Logger.Error("Failed to get permission menu by UUID", err)
		return nil, errors.New("failed to get permission menu by UUID")
	}
	return permissionMenu, nil
}

// UpdatePermissionMenu 更新权限菜单关联信息
func (s *PermissionMenuService) UpdatePermissionMenu(ctx *app.Context, permissionMenu *model.PermissionMenu) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	permissionMenu.UpdatedAt = now
	err := ctx.DB.Where("uuid = ?", permissionMenu.Uuid).Updates(permissionMenu).Error
	if err != nil {
		ctx.Logger.Error("Failed to update permission menu", err)
		return errors.New("failed to update permission menu")
	}

	return nil
}

// DeletePermissionMenu 删除权限菜单关联
func (s *PermissionMenuService) DeletePermissionMenu(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.PermissionMenu{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete permission menu", err)
		return errors.New("failed to delete permission menu")
	}

	return nil
}

// GetPermissionMenuList 根据查询参数获取权限菜单关联列表
func (s *PermissionMenuService) GetPermissionMenuList(ctx *app.Context, params *model.ReqPermissionMenuQueryParam) (*model.PagedResponse, error) {
	var (
		permissionMenus []*model.PermissionMenu
		total           int64
	)

	db := ctx.DB.Model(&model.PermissionMenu{})

	if params.PermissionUuid != "" {
		db = db.Where("permission_uuid = ?", params.PermissionUuid)
	}

	if params.MenuUuid != "" {
		db = db.Where("menu_uuid = ?", params.MenuUuid)
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get permission menu count", err)
		return nil, errors.New("failed to get permission menu count")
	}

	err = db.Offset(params.GetOffset()).Limit(params.PageSize).Find(&permissionMenus).Error
	if err != nil {
		ctx.Logger.Error("Failed to get permission menu list", err)
		return nil, errors.New("failed to get permission menu list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  permissionMenus,
	}, nil
}

// 根据菜单 UUID 获取权限菜单关联列表
func (s *PermissionMenuService) GetPermissionMenuListByMenuUUID(ctx *app.Context, menuUUID string) ([]*model.Permission, error) {
	var permissionMenus []*model.PermissionMenu
	err := ctx.DB.Where("menu_uuid = ?", menuUUID).Find(&permissionMenus).Error
	if err != nil {
		ctx.Logger.Error("Failed to get permission menu list by menu UUID", err)
		return nil, errors.New("failed to get permission menu list by menu UUID")
	}

	uuids := make([]string, 0)
	for _, permissionMenu := range permissionMenus {
		uuids = append(uuids, permissionMenu.PermissionUuid)
	}

	var permissions []*model.Permission
	err = ctx.DB.Where("uuid IN (?)", uuids).Find(&permissions).Error
	if err != nil {
		ctx.Logger.Error("Failed to get permission list by UUIDs", err)
		return nil, errors.New("failed to get permission list by UUIDs")
	}

	return permissions, nil
}

// 根据权限 UUID 获取权限菜单关联列表
func (s *PermissionMenuService) GetPermissionMenuListByPermissionUUID(ctx *app.Context, permissionUUID string) ([]*model.PermissionMenu, error) {

	var permissionMenus []*model.PermissionMenu
	err := ctx.DB.Where("permission_uuid = ?", permissionUUID).Find(&permissionMenus).Error
	if err != nil {
		ctx.Logger.Error("Failed to get permission menu list by permission UUID", err)
		return nil, errors.New("failed to get permission menu list by permission UUID")
	}

	return permissionMenus, nil
}
