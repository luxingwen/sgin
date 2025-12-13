package middleware

import "github.com/luxingwen/sgin/pkg/app"

// 用户权限中间件

func UserPermission() app.HandlerFunc {
	return func(c *app.Context) {

		// 获取用户id
		userId := c.GetString("user_id")

		// 获取api path
		// 获取api method

		_ = userId
	}
}
