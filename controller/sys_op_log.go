package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type SysOpLogController struct {
	SysOpLogService *service.SysOpLogService
}

// @Summary 创建操作日志
// @Description 创建操作日志
// @Tags 操作日志
// @Accept  json
// @Produce  json
// @Param param body model.SysOpLog true "操作日志参数"
// @Success 200 {object} model.SysOpLog
// @Router /api/v1/sysoplog/create [post]
func (s *SysOpLogController) CreateSysOpLog(ctx *app.Context) {
	var param model.SysOpLog
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := s.SysOpLogService.CreateSysOpLog(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 更新操作日志
// @Description 更新操作日志
// @Tags 操作日志
// @Accept  json
// @Produce  json
// @Param param body model.SysOpLog true "操作日志参数"
// @Success 200 {object} model.SysOpLog
// @Router /api/v1/sysoplog/update [post]
func (s *SysOpLogController) UpdateSysOpLog(ctx *app.Context) {
	var param model.SysOpLog
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := s.SysOpLogService.UpdateSysOpLog(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 删除操作日志
// @Description 删除操作日志
// @Tags 操作日志
// @Accept  json
// @Produce  json
// @Param param body model.ReqIdParam true "操作日志ID"
// @Success 200 {string} string "ok"
// @Router /api/v1/sysoplog/delete [post]
func (s *SysOpLogController) DeleteSysOpLog(ctx *app.Context) {
	var param model.ReqIdParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := s.SysOpLogService.DeleteSysOpLog(ctx, param.Id); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// @Summary 获取操作日志信息
// @Description 获取操作日志信息
// @Tags 操作日志
// @Accept  json
// @Produce  json
// @Param param body model.ReqIdParam true "操作日志ID"
// @Success 200 {object} model.SysOpLog
// @Router /api/v1/sysoplog/info [post]
func (s *SysOpLogController) GetSysOpLogInfo(ctx *app.Context) {
	var param model.ReqIdParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	log, err := s.SysOpLogService.GetSysOpLogByID(ctx, param.Id)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(log)
}

// @Summary 获取操作日志列表
// @Description 获取操作日志列表
// @Tags 操作日志
// @Accept  json
// @Produce  json
// @Param param body model.ReqOpLogQueryParam true "查询参数"
// @Success 200 {object} model.PagedResponse
// @Router /api/v1/sysoplog/list [post]
func (s *SysOpLogController) GetSysOpLogList(ctx *app.Context) {
	param := &model.ReqOpLogQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	logs, err := s.SysOpLogService.GetSysOpLogList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(logs)
}
