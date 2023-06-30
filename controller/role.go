package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
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
// @Success 200 {object} model.PagedResponse
// @Router /role/list [post]
func (c *RoleController) GetRoleList(ctx *app.Context) {
	param := &model.ReqRoleQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	roles, err := c.RoleService.GetRoleList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
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
// @Success 200 {object} model.Role
// @Router /role/create [post]
func (c *RoleController) CreateRole(ctx *app.Context) {
	param := &model.Role{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.RoleService.CreateRole(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

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
// @Router /role/update [post]
func (c *RoleController) UpdateRole(ctx *app.Context) {
	param := &model.Role{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.RoleService.UpdateRole(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

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
// @Router /role/delete [post]
func (c *RoleController) DeleteRole(ctx *app.Context) {
	param := &model.Role{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.RoleService.DeleteRole(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("删除角色成功")
}
