package service

import (
	"sgin/model"
	"sgin/pkg/app"
	"time"

	"github.com/google/uuid"
)

type ServerService struct {
}

func NewServerService() *ServerService {
	return &ServerService{}
}

// 创建服务
func (s *ServerService) CreateServer(ctx *app.Context, param *model.Server) error {
	param.UUID = uuid.New().String()
	param.CreateAt = time.Now()
	param.UpdateAt = time.Now()
	return ctx.DB.Create(param).Error
}

// 更新服务
func (s *ServerService) UpdateServer(ctx *app.Context, param *model.Server) error {
	param.UpdateAt = time.Now()
	return ctx.DB.Model(param).Where("uuid = ?", param.UUID).Updates(param).Error
}

// 删除服务
func (s *ServerService) DeleteServer(ctx *app.Context, uuid string) error {
	return ctx.DB.Where("uuid = ?", uuid).Delete(&model.Server{}).Error
}

// 获取服务信息
func (s *ServerService) GetServerInfo(ctx *app.Context, uuid string) (r *model.Server, err error) {
	r = &model.Server{}
	err = ctx.DB.Where("uuid = ?", uuid).First(r).Error
	return
}

// 获取服务列表
func (s *ServerService) GetServerList(ctx *app.Context, param model.ReqServerQueryParam) (r *model.PagedResponse, err error) {
	var (
		serverList []*model.Server
		total      int64
	)

	db := ctx.DB.Model(&model.Server{})

	if param.Name != "" {
		db = db.Where("name like ?", "%"+param.Name+"%")
	}

	if err = db.Offset(param.GetOffset()).Limit(param.PageSize).Find(&serverList).Error; err != nil {
		return
	}
	if err = db.Count(&total).Error; err != nil {

		return
	}

	r = &model.PagedResponse{
		Total:    total,
		Data:     serverList,
		Current:  param.Current,
		PageSize: param.PageSize,
	}
	return
}
