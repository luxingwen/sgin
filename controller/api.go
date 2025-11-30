package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

type APIController struct {
	APIService *service.APIService
}

// @Summary 创建API
// @Description 创建API
// @Tags API
// @Accept  json
// @Produce  json
// @Param param body model.API true "API参数"
// @Success 200 {object} model.API
// @Router /api/v1/api/create [post]
func (a *APIController) CreateAPI(ctx *app.Context) {
	var param model.API
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind create api params failed")
		return
	}
	if err := a.APIService.CreateAPI(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "create api failed")
		return
	}
	ctx.Logger.Infow("api created",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"api_uuid", param.UUID,
		"api_name", param.Name,
	)
	ctx.JSONSuccess(param)
}

// @Summary 更新API
// @Description 更新API
// @Tags API
// @Accept  json
// @Produce  json
// @Param param body model.API true "API参数"
// @Success 200 {object} model.API
// @Router /api/v1/api/update [post]
func (a *APIController) UpdateAPI(ctx *app.Context) {
	var param model.API
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind update api params failed")
		return
	}
	if err := a.APIService.UpdateAPI(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "update api failed")
		return
	}
	ctx.Logger.Infow("api updated",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"api_uuid", param.UUID,
		"api_name", param.Name,
	)
	ctx.JSONSuccess(param)
}

// @Summary 删除API
// @Description 删除API
// @Tags API
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "API UUID"
// @Success 200 {string} string	"ok"
// @Router /api/v1/api/delete [post]
func (a *APIController) DeleteAPI(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete api params failed")
		return
	}
	if err := a.APIService.DeleteAPI(ctx, param.Uuid); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "delete api failed", "uuid", param.Uuid)
		return
	}
	ctx.Logger.Infow("api deleted",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"api_uuid", param.Uuid,
	)
	ctx.JSONSuccess("ok")
}

// @Summary 获取API信息
// @Description 获取API信息
// @Tags API
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "API UUID"
// @Success 200 {object} model.API
// @Router /api/v1/api/info [post]
func (a *APIController) GetAPIInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind get api info params failed")
		return
	}
	api, err := a.APIService.GetAPIByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "get api info failed", "uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess(api)
}

// @Summary 获取API列表
// @Description 获取API列表
// @Tags API
// @Accept  json
// @Produce  json
// @Param param body model.ReqAPIQueryParam true "查询参数"
// @Success 200 {object} model.PagedResponse
// @Router /api/v1/api/list [post]
func (a *APIController) GetAPIList(ctx *app.Context) {
	param := &model.ReqAPIQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list api params failed")
		return
	}

	apis, err := a.APIService.GetAPIList(ctx, param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list apis failed")
		return
	}

	ctx.JSONSuccess(apis)
}
