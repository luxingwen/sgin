package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
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
// @Success 200 {object} string "Successfully fetched user data"
// @Router /app/list [post]
func (ac *AppController) GetAppList(c *app.Context) {
	param := &model.ReqAppQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	user, err := ac.AppService.GetAppList(c, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
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
// @Success 200 {object} string "Successfully fetched user data"
// @Router /app/create [post]
func (ac *AppController) CreateApp(c *app.Context) {
	var app model.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := ac.AppService.CreateApp(c, &app)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(app)
}

// @Tags App
// @Summary 更新应用
// @Description 更新应用
// @Accept  json
// @Produce  json
// @Param params body model.App true "Update app"
// @Success 200 {object} string "Successfully update user data"
// @Router /app/update [post]
func (ac *AppController) UpdateApp(c *app.Context) {
	var app model.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := ac.AppService.UpdateApp(c, &app)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(app)
}

// @Tags App
// @Summary 删除应用
// @Description 删除应用
// @Accept  json
// @Produce  json
// @Param params body model.ReqUuidParam true "Delete app"
// @Success 200 {object} string "Successfully delete user data"
// @Router /app/delete [post]
func (ac *AppController) DeleteApp(c *app.Context) {
	var app model.ReqUuidParam
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := ac.AppService.DeleteApp(c, app.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess("删除成功")
}
