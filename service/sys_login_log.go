package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type SysLoginLogService struct {
}

func NewSysLoginLogService() *SysLoginLogService {
	return &SysLoginLogService{}
}

func (s *SysLoginLogService) CreateLoginLog(ctx *app.Context, loginLog *model.SysLoginLog) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	loginLog.CreatedAt = now

	err := ctx.DB.Create(loginLog).Error
	if err != nil {
		ctx.Logger.Error("Failed to create login log", err)
		return errors.New("failed to create login log")
	}
	return nil
}

func (s *SysLoginLogService) GetLoginLogByID(ctx *app.Context, id uint) (*model.SysLoginLog, error) {
	loginLog := &model.SysLoginLog{}
	err := ctx.DB.First(loginLog, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("login log not found")
		}
		ctx.Logger.Error("Failed to get login log by ID", err)
		return nil, errors.New("failed to get login log by ID")
	}
	return loginLog, nil
}

func (s *SysLoginLogService) UpdateLoginLog(ctx *app.Context, loginLog *model.SysLoginLog) error {
	err := ctx.DB.Model(&model.SysLoginLog{}).Where("id = ?", loginLog.ID).Updates(loginLog).Error
	if err != nil {
		ctx.Logger.Error("Failed to update login log", err)
		return errors.New("failed to update login log")
	}
	return nil
}

func (s *SysLoginLogService) DeleteLoginLog(ctx *app.Context, id uint) error {
	err := ctx.DB.Delete(&model.SysLoginLog{}, id).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete login log", err)
		return errors.New("failed to delete login log")
	}
	return nil
}

// GetLoginLogList retrieves a list of login logs based on query parameters
func (s *SysLoginLogService) GetLoginLogList(ctx *app.Context, params *model.ReqLoginLogQueryParam) (*model.PagedResponse, error) {
	var (
		loginLogs []*model.SysLoginLog
		total     int64
	)

	db := ctx.DB.Model(&model.SysLoginLog{})

	if params.Username != "" {
		db = db.Where("username LIKE ?", "%"+params.Username+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get login log count", err)
		return nil, errors.New("failed to get login log count")
	}

	err = db.Offset(params.GetOffset()).Limit(params.PageSize).Order("id DESC").Find(&loginLogs).Error
	if err != nil {
		ctx.Logger.Error("Failed to get login log list", err)
		return nil, errors.New("failed to get login log list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  loginLogs,
	}, nil
}
