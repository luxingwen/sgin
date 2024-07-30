package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type PermissionController struct {
	PermissionService *service.PermissionService
}

// CreatePermission 创建新的权限
func (p *PermissionController) CreatePermission(ctx *app.Context) {
	var param model.Permission
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PermissionService.CreatePermission(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// UpdatePermission 更新权限信息
func (p *PermissionController) UpdatePermission(ctx *app.Context) {
	var param model.Permission
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PermissionService.UpdatePermission(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// DeletePermission 删除权限
func (p *PermissionController) DeletePermission(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PermissionService.DeletePermission(ctx, param.Uuid); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// GetPermissionInfo 获取权限信息
func (p *PermissionController) GetPermissionInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	permission, err := p.PermissionService.GetPermissionByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(permission)
}

// GetPermissionList 获取权限列表
func (p *PermissionController) GetPermissionList(ctx *app.Context) {
	var param model.ReqPermissionQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	permissions, err := p.PermissionService.GetPermissionList(ctx, &param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(permissions)
}
