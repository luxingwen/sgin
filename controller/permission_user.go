package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

type UserPermissionController struct {
	UserPermissionService *service.UserPermissionService
}

// CreateUserPermission 创建新的用户权限关联
func (u *UserPermissionController) CreateUserPermission(ctx *app.Context) {
	var param model.ReqPermissionUserCreate
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind create user permission params failed")
		return
	}
	if err := u.UserPermissionService.CreateUserPermission(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "create user permission failed")
		return
	}
	ctx.JSONSuccess(param)
}

// UpdateUserPermission 更新用户权限关联信息
func (u *UserPermissionController) UpdateUserPermission(ctx *app.Context) {
	var param model.UserPermission
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind update user permission params failed")
		return
	}
	if err := u.UserPermissionService.UpdateUserPermission(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "update user permission failed")
		return
	}
	ctx.JSONSuccess(param)
}

// DeleteUserPermission 删除用户权限关联
func (u *UserPermissionController) DeleteUserPermission(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete user permission params failed")
		return
	}
	if err := u.UserPermissionService.DeleteUserPermission(ctx, param.Uuid); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "delete user permission failed", "uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess("ok")
}

// GetUserPermissionInfo 获取用户权限关联信息
func (u *UserPermissionController) GetUserPermissionInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind get user permission info params failed")
		return
	}
	userPermission, err := u.UserPermissionService.GetUserPermissionByUserUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "get user permission info failed", "user_uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess(userPermission)
}

// GetUserPermissionList 获取用户权限关联列表
func (u *UserPermissionController) GetUserPermissionList(ctx *app.Context) {
	var param model.ReqUserPermissionQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list user permissions params failed")
		return
	}
	userPermissions, err := u.UserPermissionService.GetUserPermissionList(ctx, &param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list user permissions failed")
		return
	}
	ctx.JSONSuccess(userPermissions)
}
