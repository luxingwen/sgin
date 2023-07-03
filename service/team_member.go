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

// 获取团队成员用户列表
func (s *TeamMemberService) GetTeamMemberUserList(ctx *app.Context, teamUUID string) ([]*model.User, error) {
	var teamMembers []*model.TeamMember

	err := ctx.DB.Where("team_uuid = ?", teamUUID).Find(&teamMembers).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("team member not found")
		}
		ctx.Logger.Error("Failed to get team member by UUID", err)
		return nil, errors.New("failed to get team member by UUID")
	}

	var users []*model.User
	var userIds []string
	for _, teamMember := range teamMembers {
		userIds = append(userIds, teamMember.UserUUID)
	}

	err = ctx.DB.Where("uuid in ?", userIds).Find(&users).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		ctx.Logger.Error("Failed to get user by UUID", err)
		return nil, errors.New("failed to get user by UUID")
	}

	return users, nil
}
