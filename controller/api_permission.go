package controller

import (
	"sgin/pkg/app"
	"sgin/service"
)

type ApiPermissionController struct {
	APIPermissionService *service.AppPermissionService
}

// @Summary Get API Permission List
// @Description Get API Permission List
// @Tags API Permission
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started" default(Bearer <Token>)
// @Param params body model.ReqApiPermissionParam true "查询参数"
// @Success 200 {object} model.PagedResponse
// @Router /api/v1/app/api/permissions/list [post]
func (ac *ApiPermissionController) List(c *app.Context) {
	// api 权限列表
}
