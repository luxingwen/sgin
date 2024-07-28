package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type SysOpLogService struct {
}

func NewSysOpLogService() *SysOpLogService {
	return &SysOpLogService{}
}

func (s *SysOpLogService) CreateSysOpLog(ctx *app.Context, log *model.SysOpLog) error {
	log.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	err := ctx.DB.Create(log).Error
	if err != nil {
		ctx.Logger.Error("Failed to create operation log", err)
		return errors.New("failed to create operation log")
	}
	return nil
}

func (s *SysOpLogService) GetSysOpLogByID(ctx *app.Context, id int64) (*model.SysOpLog, error) {
	log := &model.SysOpLog{}
	err := ctx.DB.Where("id = ?", id).First(log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("operation log not found")
		}
		ctx.Logger.Error("Failed to get operation log by ID", err)
		return nil, errors.New("failed to get operation log by ID")
	}
	return log, nil
}

func (s *SysOpLogService) UpdateSysOpLog(ctx *app.Context, log *model.SysOpLog) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	log.CreatedAt = now
	err := ctx.DB.Where("id = ?", log.ID).Updates(log).Error
	if err != nil {
		ctx.Logger.Error("Failed to update operation log", err)
		return errors.New("failed to update operation log")
	}
	return nil
}

func (s *SysOpLogService) DeleteSysOpLog(ctx *app.Context, id int64) error {
	err := ctx.DB.Model(&model.SysOpLog{}).Where("id = ?", id).Delete(&model.SysOpLog{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete operation log", err)
		return errors.New("failed to delete operation log")
	}
	return nil
}

// GetSysOpLogList retrieves a list of operation logs based on query parameters
func (s *SysOpLogService) GetSysOpLogList(ctx *app.Context, params *model.ReqOpLogQueryParam) (*model.PagedResponse, error) {
	var (
		logs  []*model.SysOpLog
		total int64
	)

	db := ctx.DB.Model(&model.SysOpLog{})

	if params.UserName != "" {
		//db = db.Where("user_name LIKE ?", "%"+params.UserName+"%")
	}

	if params.Path != "" {
		db = db.Where("path LIKE ?", "%"+params.Path+"%")
	}

	if params.Method != "" {
		db = db.Where("method = ?", params.Method)
	}

	if params.Status != 0 {
		db = db.Where("status = ?", params.Status)
	}

	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get operation log count", err)
		return nil, errors.New("failed to get operation log count")
	}

	err = db.Order("id DESC").Offset(params.GetOffset()).Limit(params.PageSize).Find(&logs).Error
	if err != nil {
		ctx.Logger.Error("Failed to get operation log list", err)
		return nil, errors.New("failed to get operation log list")
	}

	userUuids := make([]string, 0)
	paths := make([]string, 0)
	for _, v := range logs {
		userUuids = append(userUuids, v.UserUuid)
		paths = append(paths, v.Path)
	}

	userMap, err := NewUserService().GetUsersByUUIDs(ctx, userUuids)
	if err != nil {
		ctx.Logger.Error("Failed to get user list by UUIDs", err)
		return nil, errors.New("failed to get user list by UUIDs")
	}

	apiMap, err := NewAPIService().GetAPIByPathList(ctx, paths)
	if err != nil {
		ctx.Logger.Error("Failed to get API list by paths", err)
		return nil, errors.New("failed to get API list by paths")
	}

	res := make([]*model.SysOpLogRes, 0)
	for _, log := range logs {
		logRes := &model.SysOpLogRes{
			SysOpLog: *log,
		}
		if user, ok := userMap[log.UserUuid]; ok {
			logRes.Username = user.Nickname
		}

		if api, ok := apiMap[log.Path]; ok {
			logRes.Module = api.Module
			logRes.Name = api.Name
		}
		res = append(res, logRes)
	}

	return &model.PagedResponse{
		Total: total,
		Data:  res,
	}, nil
}
