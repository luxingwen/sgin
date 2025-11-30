package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

type RoleController struct {
	RoleService *service.RoleService
}

// 查询角色列表
// @Summary 查询角色列表
// @Tags 角色
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqRoleQueryParam false "查询参数"
// @Success 200 {object} model.RoleQueryResponse
// @Router /api/v1/role/list [post]
func (c *RoleController) GetRoleList(ctx *app.Context) {
	param := &model.ReqRoleQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list role params failed")
		return
	}

	roles, err := c.RoleService.GetRoleList(ctx, param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list roles failed")
		return
	}

	ctx.JSONSuccess(roles)
}

// 创建角色
// @Summary 创建角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Role true "角色信息"
// @Success 200 {object} model.RoleInfoResponse
// @Router /api/v1/role/create [post]
func (c *RoleController) CreateRole(ctx *app.Context) {
	param := &model.Role{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind create role params failed")
		return
	}

	err := c.RoleService.CreateRole(ctx, param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "create role failed")
		return
	}
	ctx.Logger.Infow("role created",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"role_uuid", param.Uuid,
		"role_name", param.Name,
	)
	ctx.JSONSuccess(param)
}

// 更新角色
// @Summary 更新角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Role true "角色信息"
// @Success 200 {object} app.Response "Successfully updated role"
// @Router /api/v1/role/update [post]
func (c *RoleController) UpdateRole(ctx *app.Context) {
	param := &model.Role{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind update role params failed")
		return
	}

	err := c.RoleService.UpdateRole(ctx, param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "update role failed")
		return
	}
	ctx.Logger.Infow("role updated",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"role_uuid", param.Uuid,
		"role_name", param.Name,
	)
	ctx.JSONSuccess("更新角色成功")
}

// 删除角色
// @Summary 删除角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Role true "角色ID"
// @Success 200 {object} app.Response "Successfully deleted role"
// @Router /api/v1/role/delete [post]
func (c *RoleController) DeleteRole(ctx *app.Context) {
	param := &model.Role{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete role params failed")
		return
	}

	err := c.RoleService.DeleteRole(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "delete role failed", "uuid", param.Uuid)
		return
	}
	ctx.Logger.Infow("role deleted",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"role_uuid", param.Uuid,
	)
	ctx.JSONSuccess("删除角色成功")
}
