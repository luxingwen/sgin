package middleware

import (
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/pkg/utils"
)

// 登录中间件
func LoginCheck() app.HandlerFunc {
	return func(c *app.Context) {

		// 获取token
		token := c.GetHeader("X-Token")
		if token == "" {
			// 兼容 Authorization: Bearer <token>
			auth := c.GetHeader("Authorization")
			const prefix = "Bearer "
			if len(auth) > len(prefix) && auth[:len(prefix)] == prefix {
				token = auth[len(prefix):]
			}
		}

		if token == "" {
			c.JSONErrLog(ecode.Unauthorized("missing token"), "missing token",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
			)
			c.Abort()
			return
		}

		// 根据token获取用户信息
		userId, err := utils.ParseTokenGetUserID(token)
		if err != nil {
			// 避免在日志中输出完整 token（敏感），但记录一个 masked 片段方便排查
			snippet := token
			if len(snippet) > 20 {
				snippet = snippet[:6] + "..." + snippet[len(snippet)-6:]
			}
			c.JSONErrLog(ecode.Unauthorized("invalid token"), "parse token failed",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"cause", err.Error(),
				"token_snippet", snippet,
			)
			c.Abort()
			return
		}

		// 将用户信息放入上下文
		c.Set("user_id", userId)
	}
}
