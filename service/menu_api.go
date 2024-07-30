package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuAPIService struct {
}

func NewMenuAPIService() *MenuAPIService {
	return &MenuAPIService{}
}

// CreateMenuAPI 创建新的菜单API关联
func (s *MenuAPIService) CreateMenuAPI(ctx *app.Context, menuAPI *model.ReqMenuAPICreate) error {
	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// 先删除已有的菜单API关联
		err := tx.Where("menu_uuid = ?", menuAPI.MenuUUID).Delete(&model.MenuAPI{}).Error
		if err != nil {
			ctx.Logger.Error("Failed to delete menu API by menu UUID", err)
			return errors.New("failed to delete menu API by menu UUID")
		}

		// 创建新的菜单API关联
		rlist := make([]*model.MenuAPI, 0)
		for _, apiUUID := range menuAPI.APIUUIDs {
			menuAPI := &model.MenuAPI{
				Uuid:      uuid.New().String(),
				MenuUUID:  menuAPI.MenuUUID,
				APIUUID:   apiUUID,
				CreatedAt: now,
				UpdatedAt: now,
			}
			rlist = append(rlist, menuAPI)
		}

		err = tx.Create(&rlist).Error
		if err != nil {
			ctx.Logger.Error("Failed to create menu API", err)
			return errors.New("failed to create menu API")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// GetMenuAPIByUUID 根据UUID获取菜单API关联
func (s *MenuAPIService) GetMenuAPIByUUID(ctx *app.Context, uuid string) (*model.MenuAPI, error) {
	menuAPI := &model.MenuAPI{}
	err := ctx.DB.Where("uuid = ?", uuid).First(menuAPI).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("menu API not found")
		}
		ctx.Logger.Error("Failed to get menu API by UUID", err)
		return nil, errors.New("failed to get menu API by UUID")
	}
	return menuAPI, nil
}

// UpdateMenuAPI 更新菜单API关联信息
func (s *MenuAPIService) UpdateMenuAPI(ctx *app.Context, menuAPI *model.MenuAPI) error {
	now := time.Now()
	menuAPI.UpdatedAt = now
	err := ctx.DB.Where("uuid = ?", menuAPI.Uuid).Updates(menuAPI).Error
	if err != nil {
		ctx.Logger.Error("Failed to update menu API", err)
		return errors.New("failed to update menu API")
	}

	return nil
}

// DeleteMenuAPI 删除菜单API关联
func (s *MenuAPIService) DeleteMenuAPI(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.MenuAPI{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete menu API", err)
		return errors.New("failed to delete menu API")
	}

	return nil
}

// GetMenuAPIList 获取菜单API关联列表
func (s *MenuAPIService) GetMenuAPIList(ctx *app.Context, params *model.ReqMenuAPIQueryParam) (*model.PagedResponse, error) {
	var (
		menuAPIs []*model.MenuAPI
		total    int64
	)

	db := ctx.DB.Model(&model.MenuAPI{})

	if params.MenuUUID != "" {
		db = db.Where("menu_uuid = ?", params.MenuUUID)
	}

	if params.APIUUID != "" {
		db = db.Where("api_uuid = ?", params.APIUUID)
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get menu API count", err)
		return nil, errors.New("failed to get menu API count")
	}

	err = db.Offset(params.GetOffset()).Limit(params.PageSize).Find(&menuAPIs).Error
	if err != nil {
		ctx.Logger.Error("Failed to get menu API list", err)
		return nil, errors.New("failed to get menu API list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  menuAPIs,
	}, nil
}

// GetMenuAPIListByMenuUUID 根据菜单UUID获取菜单API关联列表
func (s *MenuAPIService) GetMenuAPIListByMenuUUID(ctx *app.Context, menuUUID string) ([]*model.API, error) {
	var menuAPIs []*model.MenuAPI
	err := ctx.DB.Where("menu_uuid = ?", menuUUID).Find(&menuAPIs).Error
	if err != nil {
		ctx.Logger.Error("Failed to get menu API list by menu UUID", err)
		return nil, errors.New("failed to get menu API list by menu UUID")
	}

	// 获取API列表
	apiUUIDs := make([]string, 0)
	for _, menuAPI := range menuAPIs {
		apiUUIDs = append(apiUUIDs, menuAPI.APIUUID)
	}

	var apis []*model.API
	err = ctx.DB.Where("uuid IN (?)", apiUUIDs).Find(&apis).Error
	if err != nil {
		ctx.Logger.Error("Failed to get API list by UUIDs", err)
		return nil, errors.New("failed to get API list by UUIDs")
	}

	return apis, nil
}

// GetMenuAPIListByAPIUUID 根据API UUID获取菜单API关联列表
func (s *MenuAPIService) GetMenuAPIListByAPIUUID(ctx *app.Context, apiUUID string) ([]*model.MenuAPI, error) {
	var menuAPIs []*model.MenuAPI
	err := ctx.DB.Where("api_uuid = ?", apiUUID).Find(&menuAPIs).Error
	if err != nil {
		ctx.Logger.Error("Failed to get menu API list by API UUID", err)
		return nil, errors.New("failed to get menu API list by API UUID")
	}

	return menuAPIs, nil
}
