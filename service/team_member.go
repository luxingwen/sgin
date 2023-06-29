package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type TeamMemberService struct {
}

func NewTeamMemberService() *TeamMemberService {
	return &TeamMemberService{}
}

func (s *TeamMemberService) CreateTeamMember(ctx *app.Context, teamMember *model.TeamMember) error {
	teamMember.CreatedAt = time.Now()
	teamMember.UpdatedAt = teamMember.CreatedAt

	err := ctx.DB.Create(teamMember).Error
	if err != nil {
		ctx.Logger.Error("Failed to create team member", err)
		return errors.New("failed to create team member")
	}
	return nil
}

func (s *TeamMemberService) GetTeamMemberByUUID(ctx *app.Context, uuid string) (*model.TeamMember, error) {
	teamMember := &model.TeamMember{}
	err := ctx.DB.Where("uuid = ?", uuid).First(teamMember).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("team member not found")
		}
		ctx.Logger.Error("Failed to get team member by UUID", err)
		return nil, errors.New("failed to get team member by UUID")
	}
	return teamMember, nil
}

func (s *TeamMemberService) UpdateTeamMember(ctx *app.Context, teamMember *model.TeamMember) error {
	teamMember.UpdatedAt = time.Now()
	err := ctx.DB.Save(teamMember).Error
	if err != nil {
		ctx.Logger.Error("Failed to update team member", err)
		return errors.New("failed to update team member")
	}

	return nil
}

func (s *TeamMemberService) DeleteTeamMember(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.TeamMember{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete team member", err)
		return errors.New("failed to delete team member")
	}

	return nil
}
