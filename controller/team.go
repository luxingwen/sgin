package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
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
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := t.TeamService.CreateTeam(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
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
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := t.TeamService.UpdateTeam(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
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
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := t.TeamService.DeleteTeam(ctx, param.Uuid); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
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
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	team, err := t.TeamService.GetTeamByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
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
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	teams, err := t.TeamService.GetTeamList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(teams)
}
