package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type ServerController struct {
	ServerService *service.ServerService
}

// @Summary 创建服务
// @Description 创建服务
// @Tags 服务
// @Accept  json
// @Produce  json
// @Param param body model.Server true "服务参数"
// @Success 200 {object} model.ServerInfoResponse
// @Router /api/v1/server/create [post]
func (s *ServerController) CreateServer(ctx *app.Context) {
	var param model.Server
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := s.ServerService.CreateServer(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 更新服务
// @Description 更新服务
// @Tags 服务
// @Accept  json
// @Produce  json
// @Param param body model.Server true "服务参数"
// @Success 200 {object} model.ServerInfoResponse
// @Router /api/v1/server/update [post]
func (s *ServerController) UpdateServer(ctx *app.Context) {
	var param model.Server
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := s.ServerService.UpdateServer(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 删除服务
// @Description 删除服务
// @Tags 服务
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "服务uuid"
// @Success 200 {string} string	"ok"
// @Router /api/v1/server/delete [post]
func (s *ServerController) DeleteServer(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := s.ServerService.DeleteServer(ctx, param.Uuid); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// @Summary 获取服务列表
// @Description 获取服务列表
// @Tags 服务
// @Accept  json
// @Produce  json
// @Param param body model.ReqServerQueryParam false "服务参数"
// @Success 200 {object} model.ServerQueryResponse
// @Router /api/v1/server/list [post]
func (s *ServerController) GetServerList(ctx *app.Context) {
	var param model.ReqServerQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	r, err := s.ServerService.GetServerList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(r)
}

// @Summary 获取服务信息
// @Description 获取服务信息
// @Tags 服务
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "服务uuid"
// @Success 200 {object} model.ServerInfoResponse
// @Router /api/v1/server/info [post]
func (s *ServerController) GetServerInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	r, err := s.ServerService.GetServerInfo(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(r)
}
