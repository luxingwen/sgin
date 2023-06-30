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
