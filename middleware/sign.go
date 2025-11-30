package middleware

import (
	"bytes"
	"io"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/pkg/utils"
	"sgin/service"
)

const maxSignBodyBytes = 1 << 20 // 1MB

// 签名校验中间件
func Signature() app.HandlerFunc {
	return func(c *app.Context) {
		signature := c.GetHeader("X-Signature")
		if signature == "" {
			c.Next()
			return
		}

		appId := c.GetHeader("X-App-Id")
		if appId == "" {
			c.JSONErrLog(ecode.Forbidden("X-App-Id is empty"), "signature missing app id",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
			)
			c.Abort()
			return
		}

		appInfo, err := service.NewAppService().GetAppByUUID(c, appId)
		if err != nil {
			c.JSONErrLog(ecode.Forbidden("invalid app"), "get app by uuid failed",
				"app_id", appId,
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"cause", err.Error(),
			)
			c.Abort()
			return
		}

		body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxSignBodyBytes))
		if err != nil {
			c.JSONErrLog(ecode.BadRequest("read request body failed"), "read request body failed",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"app_id", appId,
				"cause", err.Error(),
			)
			return
		}

		// 将 body 内容写回
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		serverSign := utils.SignBody(body, []byte(appInfo.ApiKey))

		if serverSign != signature {
			c.JSONErrLog(ecode.Forbidden("signature is invalid"), "signature is invalid",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"app_id", appId,
			)
			c.Abort()
			return
		}

		// 将 app_id 写入上下文，供限流等中间件使用
		c.Set("app_id", appId)

		c.Next()
	}
}
