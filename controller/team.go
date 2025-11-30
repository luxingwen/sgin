package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

type TeamController struct {
	TeamService *service.TeamService
}

// @Summary 创建团队
// @Description 创建团队
// @Tags 团队
// @Accept  json
// @Produce  json
// @Param param body model.Team true "团队参数"
// @Success 200 {object} model.TeamInfoResponse
// @Router /api/v1/team/create [post]
func (t *TeamController) CreateTeam(ctx *app.Context) {
	var param model.Team
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind create team params failed")
		return
	}
	if err := t.TeamService.CreateTeam(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "create team failed")
		return
	}
	ctx.Logger.Infow("team created",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"team_uuid", param.UUID,
		"team_name", param.Name,
	)
	ctx.JSONSuccess(param)
}

// @Summary 更新团队
// @Description 更新团队
// @Tags 团队
// @Accept  json
// @Produce  json
// @Param param body model.Team true "团队参数"
// @Success 200 {object} model.TeamInfoResponse
// @Router /api/v1/team/update [post]
func (t *TeamController) UpdateTeam(ctx *app.Context) {
	var param model.Team
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind update team params failed")
		return
	}
	if err := t.TeamService.UpdateTeam(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "update team failed")
		return
	}
	ctx.Logger.Infow("team updated",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"team_uuid", param.UUID,
		"team_name", param.Name,
	)
	ctx.JSONSuccess(param)
}

// @Summary 删除团队
// @Description 删除团队
// @Tags 团队
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "团队UUID"
// @Success 200 {string} string	"ok"
// @Router /api/v1/team/delete [post]
func (t *TeamController) DeleteTeam(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete team params failed")
		return
	}
	if err := t.TeamService.DeleteTeam(ctx, param.Uuid); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "delete team failed", "uuid", param.Uuid)
		return
	}
	ctx.Logger.Infow("team deleted",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"team_uuid", param.Uuid,
	)
	ctx.JSONSuccess("ok")
}

// @Summary 获取团队信息
// @Description 获取团队信息
// @Tags 团队
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "团队UUID"
// @Success 200 {object} model.TeamInfoResponse
// @Router /api/v1/team/info [post]
func (t *TeamController) GetTeamInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind get team info params failed")
		return
	}
	team, err := t.TeamService.GetTeamByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "get team info failed", "uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess(team)
}

// @Summary 获取团队列表
// @Description 获取团队列表
// @Tags 团队
// @Accept  json
// @Produce  json
// @Param param body model.ReqTeamQueryParam true "查询参数"
// @Success 200 {object} model.TeamQueryResponse
// @Router /api/v1/team/list [post]
func (t *TeamController) GetTeamList(ctx *app.Context) {
	param := &model.ReqTeamQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list teams params failed")
		return
	}

	teams, err := t.TeamService.GetTeamList(ctx, param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list teams failed")
		return
	}

	ctx.JSONSuccess(teams)
}
