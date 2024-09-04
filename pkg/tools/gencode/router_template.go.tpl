
func Init{{.StructName}}Router(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		{{.LowerStructName}}Controller := &controller.{{.StructName}}Controller{
			{{.StructName}}Service: &service.{{.StructName}}Service{},
		}
		v1.POST("/{{.LowerStructName}}/create", {{.LowerStructName}}Controller.Create{{.StructName}})
		v1.POST("/{{.LowerStructName}}/update", {{.LowerStructName}}Controller.Update{{.StructName}})
		v1.POST("/{{.LowerStructName}}/delete", {{.LowerStructName}}Controller.Delete{{.StructName}})
		v1.POST("/{{.LowerStructName}}/info", {{.LowerStructName}}Controller.Get{{.StructName}}Info)
		v1.POST("/{{.LowerStructName}}/list", {{.LowerStructName}}Controller.Get{{.StructName}}List)
	}
}