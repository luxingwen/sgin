package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

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
	err := ctx.DB.Save(menu).Error
	if err != nil {
		ctx.Logger.Error("Failed to update menu", err)
		return errors.New("failed to update menu")
	}

	return nil
}

func (s *MenuService) DeleteMenu(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.Menu{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete menu", err)
		return errors.New("failed to delete menu")
	}

	return nil
}
