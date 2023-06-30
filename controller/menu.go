package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type MenuController struct {
	MenuService *service.MenuService
}

// 查询菜单列表
// @Summary 查询菜单列表
// @Tags 菜单
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqMenuQueryParam false "查询参数"
// @Success 200 {object} model.PagedResponse
// @Router /menu/list [post]
func (c *MenuController) GetMenuList(ctx *app.Context) {
	param := &model.ReqMenuQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	menus, err := c.MenuService.GetMenuList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(menus)
}
