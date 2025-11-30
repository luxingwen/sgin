package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

type AppController struct {
	AppService *service.AppService
}

// @Tags App
// @Summary 获取应用列表
// @Description 获取应用列表
// @Accept  json
// @Produce  json
// @Param params body model.ReqAppQueryParam false "查询参数"
// @Success 200 {object} model.AppQueryResponse
// @Router /api/v1/app/list [post]
func (ac *AppController) GetAppList(c *app.Context) {
	param := &model.ReqAppQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind list app params failed")
		return
	}

	user, err := ac.AppService.GetAppList(c, param)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "list apps failed")
		return
	}
	c.JSONSuccess(user)
}

// @Tags App
// @Summary 创建应用
// @Description 创建应用
// @Accept  json
// @Produce  json
// @Param params body model.App true "Create app"
// @Success 200 {object} model.AppInfoResponse
// @Router /api/v1/app/create [post]
func (ac *AppController) CreateApp(c *app.Context) {
	var app model.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind create app params failed")
		return
	}

	err := ac.AppService.CreateApp(c, &app)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "create app failed", "name", app.Name)
		return
	}
	c.Logger.Infow("app created",
		"path", c.FullPath(),
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
		"app_uuid", app.UUID,
		"app_name", app.Name,
	)
	c.JSONSuccess(app)
}

// @Tags App
// @Summary 更新应用
// @Description 更新应用
// @Accept  json
// @Produce  json
// @Param params body model.App true "Update app"
// @Success 200 {object} model.AppInfoResponse
// @Router /api/v1/app/update [post]
func (ac *AppController) UpdateApp(c *app.Context) {
	var app model.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind update app params failed")
		return
	}

	err := ac.AppService.UpdateApp(c, &app)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "update app failed", "uuid", app.UUID)
		return
	}
	c.Logger.Infow("app updated",
		"path", c.FullPath(),
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
		"app_uuid", app.UUID,
		"app_name", app.Name,
	)
	c.JSONSuccess(app)
}

// @Tags App
// @Summary 删除应用
// @Description 删除应用
// @Accept  json
// @Produce  json
// @Param params body model.ReqUuidParam true "Delete app"
// @Success 200 {object} app.Response "Successfully delete user data"
// @Router /api/v1/app/delete [post]
func (ac *AppController) DeleteApp(c *app.Context) {
	var app model.ReqUuidParam
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete app params failed")
		return
	}

	err := ac.AppService.DeleteApp(c, app.Uuid)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "delete app failed", "uuid", app.Uuid)
		return
	}
	c.Logger.Infow("app deleted",
		"path", c.FullPath(),
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
		"app_uuid", app.Uuid,
	)
	c.JSONSuccess("删除成功")
}
