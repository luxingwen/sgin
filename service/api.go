package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APIService struct {
}

func NewAPIService() *APIService {
	return &APIService{}
}

func (s *APIService) CreateAPI(ctx *app.Context, api *model.API) error {
	now := time.Now()
	api.CreatedAt = now
	api.UpdatedAt = now
	api.UUID = uuid.New().String()

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
	now := time.Now()
	api.UpdatedAt = now
	err := ctx.DB.Where("uuid = ?", api.UUID).Updates(api).Error
	if err != nil {
		ctx.Logger.Error("Failed to update API", err)
		return errors.New("failed to update API")
	}

	return nil
}

func (s *APIService) DeleteAPI(ctx *app.Context, uuid string) error {
	err := ctx.DB.Model(&model.API{}).Where("uuid = ?", uuid).Update("status", 2).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete API", err)
		return errors.New("failed to delete API")
	}

	return nil
}

func (s *APIService) GetAPIList(ctx *app.Context, params *model.ReqAPIQueryParam) (*model.PagedResponse, error) {
	var (
		apis  []*model.API
		total int64
	)

	db := ctx.DB.Model(&model.API{})

	if params.Module != "" {
		db = db.Where("module LIKE ?", "%"+params.Module+"%")
	}
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.Status != 0 {
		db = db.Where("status = ?", params.Status)
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get API count", err)
		return nil, errors.New("failed to get API count")
	}

	err = db.Offset(params.GetOffset()).Limit(params.PageSize).Find(&apis).Error
	if err != nil {
		ctx.Logger.Error("Failed to get API list", err)
		return nil, errors.New("failed to get API list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  apis,
	}, nil
}

// 根据path列表 获取API
func (s *APIService) GetAPIByPathList(ctx *app.Context, pathList []string) (map[string]*model.API, error) {
	var (
		apis []*model.API
	)
	apiMap := make(map[string]*model.API)

	err := ctx.DB.Model(&model.API{}).Where("path IN ?", pathList).Find(&apis).Error
	if err != nil {
		ctx.Logger.Error("Failed to get API by path list", err)
		return nil, errors.New("failed to get API by path list")
	}

	for _, api := range apis {
		apiMap[api.Path] = api
	}

	return apiMap, nil
}
