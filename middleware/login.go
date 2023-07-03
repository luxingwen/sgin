package middleware

import (
	"net/http"
	"sgin/pkg/app"
	"sgin/pkg/utils"
)

// 登录中间件
func LoginCheck() app.HandlerFunc {
	return func(c *app.Context) {

		// 获取token
		token := c.GetHeader("X-Token")

		// 根据token获取用户信息
		userId, err := utils.ParseTokenGetUserID(token)
		if err != nil {
			c.JSONError(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// 将用户信息放入上下文
		c.Set("user_id", userId)
	}
}
