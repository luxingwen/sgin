package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

type PermissionController struct {
	PermissionService *service.PermissionService
}

// CreatePermission 创建新的权限
func (p *PermissionController) CreatePermission(ctx *app.Context) {
	var param model.Permission
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind create permission params failed")
		return
	}
	if err := p.PermissionService.CreatePermission(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "create permission failed")
		return
	}
	ctx.JSONSuccess(param)
}

// UpdatePermission 更新权限信息
func (p *PermissionController) UpdatePermission(ctx *app.Context) {
	var param model.Permission
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind update permission params failed")
		return
	}
	if err := p.PermissionService.UpdatePermission(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "update permission failed")
		return
	}
	ctx.JSONSuccess(param)
}

// DeletePermission 删除权限
func (p *PermissionController) DeletePermission(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete permission params failed")
		return
	}
	if err := p.PermissionService.DeletePermission(ctx, param.Uuid); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "delete permission failed", "uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess("ok")
}

// GetPermissionInfo 获取权限信息
func (p *PermissionController) GetPermissionInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind get permission info params failed")
		return
	}
	permission, err := p.PermissionService.GetPermissionByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "get permission info failed", "uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess(permission)
}

// GetPermissionList 获取权限列表
func (p *PermissionController) GetPermissionList(ctx *app.Context) {
	var param model.ReqPermissionQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list permission params failed")
		return
	}
	permissions, err := p.PermissionService.GetPermissionList(ctx, &param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list permissions failed")
		return
	}
	ctx.JSONSuccess(permissions)
}
