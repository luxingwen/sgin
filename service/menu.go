package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuService struct {
}

func NewMenuService() *MenuService {
	return &MenuService{}
}

func (s *MenuService) CreateMenu(ctx *app.Context, menu *model.Menu) error {
	menu.CreatedAt = time.Now()
	menu.UpdatedAt = menu.CreatedAt
	menu.UUID = uuid.New().String()

	err := ctx.DB.Create(menu).Error
	if err != nil {
		ctx.Logger.Error("Failed to create menu", err)
		return errors.New("failed to create menu")
	}
	return nil
}

func (s *MenuService) GetMenuByUUID(ctx *app.Context, uuid string) (*model.Menu, error) {
	menu := &model.Menu{}
	err := ctx.DB.Where("uuid = ?", uuid).First(menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("menu not found")
		}
		ctx.Logger.Error("Failed to get menu by UUID", err)
		return nil, errors.New("failed to get menu by UUID")
	}
	return menu, nil
}

func (s *MenuService) UpdateMenu(ctx *app.Context, menu *model.Menu) error {
	menu.UpdatedAt = time.Now()
	err := ctx.DB.Where("uuid = ?", menu.UUID).Updates(menu).Error
	if err != nil {
		ctx.Logger.Error("Failed to update menu", err)
		return errors.New("failed to update menu")
	}

	return nil
}

func (s *MenuService) DeleteMenu(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ? OR parent_uuid = ?", uuid, uuid).Delete(&model.Menu{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete menu", err)
		return errors.New("failed to delete menu")
	}

	return nil
}

// 获取菜单列表
func (s *MenuService) GetMenuList(ctx *app.Context, params *model.ReqMenuQueryParam) (r *model.PagedResponse, err error) {
	var (
		menus []*model.Menu
		total int64
	)

	db := ctx.DB.Model(&model.Menu{})

	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get menu count", err)
		return nil, errors.New("failed to get menu count")
	}

	err = db.Find(&menus).Error
	if err != nil {
		ctx.Logger.Error("Failed to get menu list", err)
		return nil, errors.New("failed to get menu list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  menus,
	}, nil
}
