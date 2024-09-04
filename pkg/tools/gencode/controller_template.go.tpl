package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type {{.StructName}}Controller struct {
	{{.StructName}}Service *service.{{.StructName}}Service
}

// @Tags {{.ModuleName}}
// @Summary 获取{{.ModuleName}}列表
// @Description 获取{{.ModuleName}}列表
// @Accept  json
// @Produce  json
// @Param params body model.{{.QueryStructName}} false "查询参数"
// @Success 200 {object} model.{{.StructName}}QueryResponse
// @Router /api/v1/{{.LowerStructName}}/list [post]
func (pc *{{.StructName}}Controller) Get{{.StructName}}List(c *app.Context) {
	param := &model.{{.QueryStructName}}{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	{{.LowerStructName}}s, err := pc.{{.StructName}}Service.Get{{.StructName}}List(c, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess({{.LowerStructName}}s)
}

// @Tags {{.ModuleName}}
// @Summary 获取{{.ModuleName}}详情
// @Description 获取{{.ModuleName}}详情
// @Accept  json
// @Produce  json
// @Param params body model.ReqUuidParam true "Get {{.LowerStructName}} info"
// @Success 200 {object} model.{{.StructName}}InfoResponse
// @Router /api/v1/{{.LowerStructName}}/info [post]
func (pc *{{.StructName}}Controller) Get{{.StructName}}Info(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	{{.LowerStructName}}, err := pc.{{.StructName}}Service.Get{{.StructName}}ByUUID(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess({{.LowerStructName}})
}

// @Tags {{.ModuleName}}
// @Summary 创建{{.ModuleName}}
// @Description 创建{{.ModuleName}}
// @Accept  json
// @Produce  json
// @Param params body model.{{.ReqCreateStructName}} true "Create {{.LowerStructName}}"
// @Success 200 {object} model.{{.StructName}}InfoResponse
// @Router /api/v1/{{.LowerStructName}}/create [post]
func (pc *{{.StructName}}Controller) Create{{.StructName}}(c *app.Context) {
	var {{.LowerStructName}} model.{{.StructName}}
	if err := c.ShouldBindJSON(&{{.LowerStructName}}); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.{{.StructName}}Service.Create{{.StructName}}(c, &{{.LowerStructName}})
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess({{.LowerStructName}})
}

// @Tags {{.ModuleName}}
// @Summary 更新{{.ModuleName}}
// @Description 更新{{.ModuleName}}
// @Accept  json
// @Produce  json
// @Param params body model.{{.StructName}} true "Update {{.LowerStructName}}"
// @Success 200 {object} model.{{.StructName}}InfoResponse
// @Router /api/v1/{{.LowerStructName}}/update [post]
func (pc *{{.StructName}}Controller) Update{{.StructName}}(c *app.Context) {
	var {{.LowerStructName}} model.{{.StructName}}
	if err := c.ShouldBindJSON(&{{.LowerStructName}}); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.{{.StructName}}Service.Update{{.StructName}}(c, &{{.LowerStructName}})
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess({{.LowerStructName}})
}

// @Tags {{.ModuleName}}
// @Summary 删除{{.ModuleName}}
// @Description 删除{{.ModuleName}}
// @Accept  json
// @Produce  json
// @Param params body model.ReqUuidParam true "Delete {{.LowerStructName}}"
// @Success 200 {object} model.StringDataResponse "Successfully deleted {{.LowerStructName}} data"
// @Router /api/v1/{{.LowerStructName}}/delete [post]
func (pc *{{.StructName}}Controller) Delete{{.StructName}}(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.{{.StructName}}Service.Delete{{.StructName}}(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess("删除成功")
}
