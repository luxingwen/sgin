package middleware

import (
	"sgin/pkg/app"
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
			c.JSONError(403, err.Error())
			c.Abort()
			return
		}
		c.Set("app_info", appInfo)
	}
}
