package middleware

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

// API权限校验中间件
func ApiPermission() app.HandlerFunc {
	return func(c *app.Context) {

		// 获取api path
		// 获取api method
		// 获取api key

		apikey := c.GetHeader("X-Api-Key")
		apiPath := c.Request.URL.Path
		apiMethod := c.Request.Method

		// 根据apikey获取app信息

		var (
			appinfo *model.App
			err     error
		)

		appinfo0, ok := c.Get("app_info")
		if !ok {
			if apikey == "" {
				c.JSONErrLog(ecode.Forbidden("missing api key"), "missing api key",
					"trace_id", c.TraceID,
					"path", c.FullPath(),
					"method", apiMethod,
					"client_ip", c.ClientIP(),
				)
				c.Abort()
				return
			}

			appinfo, err = service.NewAppService().GetAppByApiKey(c, apikey)
			if err != nil {
				c.JSONErrLog(ecode.Forbidden("invalid api key"), "get app by apikey failed",
					"trace_id", c.TraceID,
					"path", c.FullPath(),
					"method", apiMethod,
					"client_ip", c.ClientIP(),
					"cause", err.Error(),
				)
				c.Abort()
				return
			}
		} else {
			appinfo = appinfo0.(*model.App)
		}

		// 根据app信息获取app权限

		appPermission, err := service.NewAppPermissionService().GetAPIPermissionByNamePathMethod(c, appinfo.UUID, apiPath, apiMethod)
		if err != nil || appPermission == nil {
			c.JSONErrLog(ecode.Forbidden("permission denied"), "api permission denied",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", apiMethod,
				"client_ip", c.ClientIP(),
				"app_id", appinfo.UUID,
				"cause", func() string {
					if err != nil {
						return err.Error()
					}
					return ""
				}(),
			)
			c.Abort()
			return
		}
		_ = appPermission
	}
}
