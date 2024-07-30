package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type PermissionMenuController struct {
	PermissionMenuService *service.PermissionMenuService
}

// CreatePermissionMenu 创建新的权限菜单关联
func (p *PermissionMenuController) CreatePermissionMenu(ctx *app.Context) {
	var param model.ReqPermissionMenuCreate
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PermissionMenuService.CreatePermissionMenu(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// UpdatePermissionMenu 更新权限菜单关联信息
func (p *PermissionMenuController) UpdatePermissionMenu(ctx *app.Context) {
	var param model.PermissionMenu
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PermissionMenuService.UpdatePermissionMenu(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// DeletePermissionMenu 删除权限菜单关联
func (p *PermissionMenuController) DeletePermissionMenu(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PermissionMenuService.DeletePermissionMenu(ctx, param.Uuid); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// GetPermissionMenuInfo 获取权限菜单关联信息
func (p *PermissionMenuController) GetPermissionMenuInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	permissionMenus, err := p.PermissionMenuService.GetPermissionMenuListByMenuUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(permissionMenus)
}

// GetPermissionMenuList 获取权限菜单关联列表
func (p *PermissionMenuController) GetPermissionMenuList(ctx *app.Context) {
	var param model.ReqPermissionMenuQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	permissionMenus, err := p.PermissionMenuService.GetPermissionMenuList(ctx, &param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(permissionMenus)
}

// GetPermissionMenuListByPermissionUUID 根据权限UUID获取权限菜单关联列表
func (p *PermissionMenuController) GetPermissionMenuListByPermissionUUID(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	permissionMenus, err := p.PermissionMenuService.GetPermissionMenuListByPermissionUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(permissionMenus)
}
