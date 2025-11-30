package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

type MenuAPIController struct {
	MenuAPIService *service.MenuAPIService
}

// CreateMenuAPI 创建新的菜单API关联
func (m *MenuAPIController) CreateMenuAPI(ctx *app.Context) {
	var param model.ReqMenuAPICreate
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind create menu api params failed")
		return
	}
	if err := m.MenuAPIService.CreateMenuAPI(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "create menu api failed")
		return
	}
	ctx.JSONSuccess(param)
}

// UpdateMenuAPI 更新菜单API关联信息
func (m *MenuAPIController) UpdateMenuAPI(ctx *app.Context) {
	var param model.MenuAPI
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind update menu api params failed")
		return
	}
	if err := m.MenuAPIService.UpdateMenuAPI(ctx, &param); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "update menu api failed")
		return
	}
	ctx.JSONSuccess(param)
}

// DeleteMenuAPI 删除菜单API关联
func (m *MenuAPIController) DeleteMenuAPI(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete menu api params failed")
		return
	}
	if err := m.MenuAPIService.DeleteMenuAPI(ctx, param.Uuid); err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "delete menu api failed", "uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess("ok")
}

// GetMenuAPIInfo 获取菜单API关联信息
func (m *MenuAPIController) GetMenuAPIInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind get menu api info params failed")
		return
	}
	menuAPI, err := m.MenuAPIService.GetMenuAPIByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "get menu api info failed", "uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess(menuAPI)
}

// GetMenuAPIList 获取菜单API关联列表
func (m *MenuAPIController) GetMenuAPIList(ctx *app.Context) {
	var param model.ReqMenuAPIQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list menu apis params failed")
		return
	}
	menuAPIs, err := m.MenuAPIService.GetMenuAPIList(ctx, &param)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list menu apis failed")
		return
	}
	ctx.JSONSuccess(menuAPIs)
}

// GetMenuAPIListByMenuUUID 根据菜单UUID获取菜单API关联列表
func (m *MenuAPIController) GetMenuAPIListByMenuUUID(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list by menu uuid params failed")
		return
	}
	menuAPIs, err := m.MenuAPIService.GetMenuAPIListByMenuUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list menu apis by menu uuid failed", "menu_uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess(menuAPIs)
}

// GetMenuAPIListByAPIUUID 根据API UUID获取菜单API关联列表
func (m *MenuAPIController) GetMenuAPIListByAPIUUID(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind list by api uuid params failed")
		return
	}
	menuAPIs, err := m.MenuAPIService.GetMenuAPIListByAPIUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "list menu apis by api uuid failed", "api_uuid", param.Uuid)
		return
	}
	ctx.JSONSuccess(menuAPIs)
}
