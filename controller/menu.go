package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
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
// @Router /api/v1/menu/list [post]
func (c *MenuController) GetMenuList(ctx *app.Context) {
	param := &model.ReqMenuQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list menu params failed")
		return
	}

	menus, err := c.MenuService.GetMenuList(ctx, param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list menus failed")
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
// @Router /api/v1/menu/create [post]
func (c *MenuController) CreateMenu(ctx *app.Context) {
	param := &model.Menu{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind create menu params failed")
		return
	}

	err := c.MenuService.CreateMenu(ctx, param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "create menu failed")
		return
	}
	ctx.Logger.Infow("menu created",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"menu_uuid", param.UUID,
		"menu_name", param.Name,
	)
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
// @Router /api/v1/menu/update [post]
func (c *MenuController) UpdateMenu(ctx *app.Context) {
	param := &model.Menu{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind update menu params failed")
		return
	}

	err := c.MenuService.UpdateMenu(ctx, param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "update menu failed")
		return
	}
	ctx.Logger.Infow("menu updated",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"menu_uuid", param.UUID,
		"menu_name", param.Name,
	)
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
// @Router /api/v1/menu/delete [post]
func (c *MenuController) DeleteMenu(ctx *app.Context) {
	param := &model.ReqMenuDeleteParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete menu params failed")
		return
	}

	err := c.MenuService.DeleteMenu(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "delete menu failed", "uuid", param.Uuid)
		return
	}
	ctx.Logger.Infow("menu deleted",
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
		"menu_uuid", param.Uuid,
	)
	ctx.JSONSuccess("ok")
}

func (c *MenuController) GetMenuInfo(ctx *app.Context) {
	param := &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind get menu info params failed")
		return
	}
	menu, err := c.MenuService.GetMenuByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "get menu info failed", "uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess(menu)
}
