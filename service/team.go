package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamService struct {
}

func NewTeamService() *TeamService {
	return &TeamService{}
}

func (s *TeamService) CreateTeam(ctx *app.Context, team *model.Team) error {
	team.UUID = uuid.New().String()
	team.CreatedAt = time.Now()
	team.UpdatedAt = team.CreatedAt

	err := ctx.DB.Create(team).Error
	if err != nil {
		ctx.Logger.Error("Failed to create team", err)
		return errors.New("failed to create team")
	}
	return nil
}

func (s *TeamService) GetTeamByUUID(ctx *app.Context, uuid string) (*model.Team, error) {
	team := &model.Team{}
	err := ctx.DB.Where("uuid = ?", uuid).First(team).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("team not found")
		}
		ctx.Logger.Error("Failed to get team by UUID", err)
		return nil, errors.New("failed to get team by UUID")
	}
	return team, nil
}

func (s *TeamService) UpdateTeam(ctx *app.Context, team *model.Team) error {
	team.UpdatedAt = time.Now()
	err := ctx.DB.Save(team).Error
	if err != nil {
		ctx.Logger.Error("Failed to update team", err)
		return errors.New("failed to update team")
	}

	return nil
}

func (s *TeamService) DeleteTeam(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.Team{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete team", err)
		return errors.New("failed to delete team")
	}

	return nil
}

// 获取团队列表
func (s *TeamService) GetTeamList(ctx *app.Context, param *model.ReqTeamQueryParam) (r *model.PagedResponse, err error) {
	var (
		teamList []*model.Team
		total    int64
	)

	db := ctx.DB.Model(&model.Team{})

	if param.Name != "" {
		db = db.Where("name like ?", "%"+param.Name+"%")
	}

	if err = db.Offset(param.GetOffset()).Limit(param.PageSize).Find(&teamList).Error; err != nil {
		return
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}

	r = &model.PagedResponse{
		Total:    total,
		Current:  param.Current,
		PageSize: param.PageSize,
		Data:     teamList,
	}

	return
}
