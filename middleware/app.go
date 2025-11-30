package middleware

import (
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

// app中间件，检查app的key是否有效
func AppKeyCheck() app.HandlerFunc {
	return func(c *app.Context) {

		// 获取apikey
		apikey := c.GetHeader("X-Api-Key")

		// 根据apikey获取app信息

		appInfo, err := service.NewAppService().GetAppByApiKey(c, apikey)
		if err != nil {
			c.JSONErrLog(ecode.Forbidden("invalid api key"), "get app by apikey failed",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"cause", err.Error(),
			)
			c.Abort()
			return
		}
		c.Set("app_info", appInfo)
		// 同步设置 app_id 方便后续限流等
		c.Set("app_id", appInfo.UUID)
	}
}
