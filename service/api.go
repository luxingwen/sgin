package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type APIService struct {
}

func NewAPIService() *APIService {
	return &APIService{}
}

func (s *APIService) CreateAPI(ctx *app.Context, api *model.API) error {
	api.CreatedAt = time.Now()
	api.UpdatedAt = api.CreatedAt

	err := ctx.DB.Create(api).Error
	if err != nil {
		ctx.Logger.Error("Failed to create API", err)
		return errors.New("failed to create API")
	}
	return nil
}

func (s *APIService) GetAPIByUUID(ctx *app.Context, uuid string) (*model.API, error) {
	api := &model.API{}
	err := ctx.DB.Where("uuid = ?", uuid).First(api).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("API not found")
		}
		ctx.Logger.Error("Failed to get API by UUID", err)
		return nil, errors.New("failed to get API by UUID")
	}
	return api, nil
}

func (s *APIService) UpdateAPI(ctx *app.Context, api *model.API) error {
	api.UpdatedAt = time.Now()
	err := ctx.DB.Save(api).Error
	if err != nil {
		ctx.Logger.Error("Failed to update API", err)
		return errors.New("failed to update API")
	}

	return nil
}

func (s *APIService) DeleteAPI(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.API{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete API", err)
		return errors.New("failed to delete API")
	}

	return nil
}
