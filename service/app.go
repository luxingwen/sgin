package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type AppService struct {
}

func NewAppService() *AppService {
	return &AppService{}
}

func (s *AppService) CreateApp(ctx *app.Context, app *model.App) error {
	app.CreatedAt = time.Now()
	app.UpdatedAt = app.CreatedAt

	err := ctx.DB.Create(app).Error
	if err != nil {
		ctx.Logger.Error("Failed to create app", err)
		return errors.New("failed to create app")
	}
	return nil
}

func (s *AppService) GetAppByUUID(ctx *app.Context, uuid string) (*model.App, error) {
	app := &model.App{}
	err := ctx.DB.Where("uuid = ?", uuid).First(app).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("app not found")
		}
		ctx.Logger.Error("Failed to get app by UUID", err)
		return nil, errors.New("failed to get app by UUID")
	}
	return app, nil
}

func (s *AppService) UpdateApp(ctx *app.Context, app *model.App) error {
	app.UpdatedAt = time.Now()
	err := ctx.DB.Save(app).Error
	if err != nil {
		ctx.Logger.Error("Failed to update app", err)
		return errors.New("failed to update app")
	}

	return nil
}

func (s *AppService) DeleteApp(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.App{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete app", err)
		return errors.New("failed to delete app")
	}

	return nil
}

// 查询app列表
func (s *AppService) GetAppList(ctx *app.Context, params *model.ReqAppQueryParam) (r *model.PagedResponse, err error) {

	var (
		apps  []*model.App
		total int64
	)

	query := ctx.DB.Model(&model.App{})

	if params.Name != "" {
		query = query.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.ApiKey != "" {
		query = query.Where("api_key LIKE ?", "%"+params.ApiKey+"%")
	}

	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}

	err = query.Offset(params.GetOffset()).Limit(params.PageSize).Find(&apps).Error
	if err != nil {
		ctx.Logger.Error("Failed to get app list", err)
		return nil, errors.New("failed to get app list")
	}

	err = query.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to count app list", err)
		return nil, errors.New("failed to count app list")
	}

	return &model.PagedResponse{
		Total:    total,
		Data:     apps,
		Current:  params.Current,
		PageSize: params.PageSize,
	}, nil
}
