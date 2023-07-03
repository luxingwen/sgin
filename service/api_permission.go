package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type AppPermissionService struct {
}

func NewAppPermissionService() *AppPermissionService {
	return &AppPermissionService{}
}

func (s *AppPermissionService) CreateAppPermission(ctx *app.Context, ap *model.AppPermission) error {
	ap.CreatedAt = time.Now()
	ap.UpdatedAt = ap.CreatedAt

	err := ctx.DB.Create(ap).Error
	if err != nil {
		ctx.Logger.Error("Failed to create app permission", err)
		return errors.New("failed to create app permission")
	}
	return nil
}

func (s *AppPermissionService) GetAppPermissionByUUID(ctx *app.Context, uuid string) (*model.AppPermission, error) {
	ap := &model.AppPermission{}
	err := ctx.DB.Where("uuid = ?", uuid).First(ap).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("app permission not found")
		}
		ctx.Logger.Error("Failed to get app permission by UUID", err)
		return nil, errors.New("failed to get app permission by UUID")
	}
	return ap, nil
}

func (s *AppPermissionService) UpdateAppPermission(ctx *app.Context, ap *model.AppPermission) error {
	ap.UpdatedAt = time.Now()
	err := ctx.DB.Save(ap).Error
	if err != nil {
		ctx.Logger.Error("Failed to update app permission", err)
		return errors.New("failed to update app permission")
	}

	return nil
}

func (s *AppPermissionService) DeleteAppPermission(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.AppPermission{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete app permission", err)
		return errors.New("failed to delete app permission")
	}

	return nil
}

// 获取app的api权限列表
func (s *AppPermissionService) GetAppAPIPermissions(ctx *app.Context, appUUID string) ([]*model.API, error) {
	var apis []*model.API
	err := ctx.DB.Table("apis").Joins("left join app_permissions on apis.uuid = app_permissions.api_uuid").Where("app_permissions.app_uuid = ?", appUUID).Find(&apis).Error
	if err != nil {
		ctx.Logger.Error("Failed to get app api permissions", err)
		return nil, errors.New("failed to get app api permissions")
	}

	return apis, nil
}

// 根据name ，path，method获取api 权限信息
func (s *AppPermissionService) GetAPIPermissionByNamePathMethod(ctx *app.Context, appUUID, path, method string) (*model.AppPermission, error) {
	appPermission := &model.AppPermission{}
	err := ctx.DB.Where("app_uuid = ? and path = ? and method = ?", appUUID, path, method).First(appPermission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("api permissions  not found")
		}
		ctx.Logger.Error("Failed to get api permissions by name path method", err)
		return nil, errors.New("failed to get api permissions by name path method")
	}
	return appPermission, nil
}
