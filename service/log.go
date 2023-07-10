package service

import (
	"errors"
	"sgin/model"
	"sgin/pkg/app"
	"time"
)

type LogService struct {
}

func NewLogService() *LogService {
	return &LogService{}
}

// 创建日志
func (s *LogService) CreateLog(ctx *app.Context, log *model.Log) error {
	log.CreatedAt = time.Now()
	log.UpdatedAt = log.CreatedAt

	err := ctx.DB.Create(log).Error
	if err != nil {
		ctx.Logger.Error("Failed to create log", err)
		return errors.New("failed to create log")
	}
	return nil
}

// 更新日志
func (s *LogService) UpdateLog(ctx *app.Context, log *model.Log) error {
	log.UpdatedAt = time.Now()
	err := ctx.DB.Where("uuid = ?", log.UUID).Updates(log).Error
	if err != nil {
		ctx.Logger.Error("Failed to update log", err)
		return errors.New("failed to update log")
	}

	return nil
}
