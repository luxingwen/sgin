package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
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
	teamMember.UUID = uuid.New().String()

	// 先查询是否存在相同的团队成员
	var isExistTeamMember model.TeamMember
	err := ctx.DB.Where("team_uuid = ? AND user_uuid = ?", teamMember.TeamUUID, teamMember.UserUUID).First(&isExistTeamMember).Error
	if err == nil && isExistTeamMember.Id > 0 {
		return errors.New("team member already exists")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Logger.Error("Failed to get team member by team UUID and user UUID", err)
		return errors.New("failed to get team member by team UUID and user UUID")
	}

	err = ctx.DB.Create(teamMember).Error
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
func (s *TeamMemberService) GetTeamMemberUserList(ctx *app.Context, params *model.ReqTeamMemberQueryParam) (*model.PagedResponse, error) {
	var teamMembers []*model.TeamMember
	var users []*model.User
	var userIds []string
	var total int64

	// 查找团队成员，并获取总数
	err := ctx.DB.Model(&model.TeamMember{}).Where("team_uuid = ?", params.TeamUUID).Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to count team members by team UUID", err)
		return nil, errors.New("failed to count team members by team UUID")
	}

	// 分页查询团队成员
	err = ctx.DB.Where("team_uuid = ?", params.TeamUUID).Offset(params.GetOffset()).Limit(params.PageSize).
		Find(&teamMembers).Error
	if err != nil {
		ctx.Logger.Error("Failed to get team members by team UUID", err)
		return nil, errors.New("failed to get team members by team UUID")
	}

	// 获取成员的用户UUID列表
	for _, teamMember := range teamMembers {
		userIds = append(userIds, teamMember.UserUUID)
	}

	// 如果没有成员，直接返回空
	if len(userIds) == 0 {
		return &model.PagedResponse{
			Total:    total,
			Data:     []*model.User{},
			Current:  params.Current,
			PageSize: params.PageSize,
		}, nil
	}

	// 查找用户
	err = ctx.DB.Where("uuid in ?", userIds).Find(&users).Error
	if err != nil {
		ctx.Logger.Error("Failed to get users by UUIDs", err)
		return nil, errors.New("failed to get users by UUIDs")
	}

	// 隐藏密码
	for _, user := range users {
		user.Password = ""
	}

	return &model.PagedResponse{
		Total:    total,
		Data:     users,
		Current:  params.Current,
		PageSize: params.PageSize,
	}, nil
}
