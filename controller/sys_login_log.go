package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

type SysLoginLogController struct {
	LoginLogService *service.SysLoginLogService
}

// @Summary 获取登录日志
// @Description 获取登录日志
// @Tags 登录日志
// @Accept  json
// @Produce  json
// @Param id path uint true "登录日志ID"
// @Success 200 {object} model.SysLoginLog
// @Router /api/v1/loginlog/{id} [get]
func (l *SysLoginLogController) GetLoginLog(ctx *app.Context) {
	// id, err := ctx.ParamUint("id")
	// if err != nil {
	// 	ctx.JSONError(http.StatusBadRequest, err.Error())
	// 	return
	// }

	// loginLog, err := l.LoginLogService.GetLoginLogByID(ctx, id)
	// if err != nil {
	// 	ctx.JSONError(http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// ctx.JSONSuccess(loginLog)
}

// @Summary 更新登录日志
// @Description 更新登录日志
// @Tags 登录日志
// @Accept  json
// @Produce  json
// @Param param body model.SysLoginLog true "登录日志参数"
// @Success 200 {object} model.SysLoginLog
// @Router /api/v1/loginlog/update [post]
func (l *SysLoginLogController) UpdateLoginLog(ctx *app.Context) {
	var param model.SysLoginLog
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind update login log params failed")
		return
	}
	if err := l.LoginLogService.UpdateLoginLog(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "update login log failed")
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 获取登录日志列表
// @Description 获取登录日志列表
// @Tags 登录日志
// @Accept  json
// @Produce  json
// @Param param body model.ReqLoginLogQueryParam true "查询参数"
// @Success 200 {object} model.PagedResponse
// @Router /api/v1/loginlog/list [post]
func (l *SysLoginLogController) GetLoginLogList(ctx *app.Context) {
	param := &model.ReqLoginLogQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list login logs params failed")
		return
	}

	loginLogs, err := l.LoginLogService.GetLoginLogList(ctx, param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list login logs failed")
		return
	}

	ctx.JSONSuccess(loginLogs)
}
