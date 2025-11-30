package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

type TeamMemberController struct {
	TeamMemberService *service.TeamMemberService
}

// @Summary 创建团队成员
// @Description 创建团队成员
// @Tags 团队成员
// @Accept  json
// @Produce  json
// @Param param body model.ReqTeamMemberCreateParam true "团队成员参数"
// @Success 200 {object} model.TeamMemberInfoResponse
// @Router /api/v1/team/member/create [post]
func (t *TeamMemberController) CreateTeamMember(ctx *app.Context) {
	var param model.ReqTeamMemberCreateParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind create team member params failed")
		return
	}

	teamMember := &model.TeamMember{
		TeamUUID: param.TeamUUID,
		UserUUID: param.UserUUID,
	}
	if err := t.TeamMemberService.CreateTeamMember(ctx, teamMember); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "create team member failed")
		return
	}
	ctx.JSONSuccess(teamMember)
}

// 删除团队成员
// @Summary 删除团队成员
// @Description 删除团队成员
// @Tags 团队成员
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "团队成员UUID"
// @Success 200 {object} app.Response
// @Router /api/v1/team/member/delete [post]
func (t *TeamMemberController) DeleteTeamMember(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete team member params failed")
		return
	}
	if err := t.TeamMemberService.DeleteTeamMember(ctx, param.Uuid); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "delete team member failed", "uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess("ok")
}

// 查询团队成员列表

func (t *TeamMemberController) GetTeamMemberList(ctx *app.Context) {
	var param model.ReqTeamMemberQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list team member params failed")
		return
	}
	teamMembers, err := t.TeamMemberService.GetTeamMemberUserList(ctx, &param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list team members failed")
		return
	}
	ctx.JSONSuccess(teamMembers)
}
