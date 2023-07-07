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
// @Success 200 {object} model.MenuQueryResponse
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

// 创建菜单
// @Summary 创建菜单
// @Tags 菜单
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Menu true "菜单信息"
// @Success 200 {object} model.MenuPageResponse
// @Router /menu/create [post]
func (c *MenuController) CreateMenu(ctx *app.Context) {
	param := &model.Menu{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.MenuService.CreateMenu(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(param)
}

// 更新菜单
// @Summary 更新菜单
// @Tags 菜单
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Menu true "菜单信息"
// @Success 200 {object} model.MenuPageResponse
// @Router /menu/update [post]
func (c *MenuController) UpdateMenu(ctx *app.Context) {
	param := &model.Menu{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.MenuService.UpdateMenu(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(param)
}

// 删除菜单
// @Summary 删除菜单
// @Tags 菜单
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqMenuDeleteParam true "删除参数"
// @Success 200 {object} app.Response
// @Router /menu/delete [post]
func (c *MenuController) DeleteMenu(ctx *app.Context) {
	param := &model.ReqMenuDeleteParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.MenuService.DeleteMenu(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("ok")
}
