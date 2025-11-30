package middleware

import (
	"fmt"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"strconv"
	"time"
)

// 防重放攻击中间件
func NonceHandler() app.HandlerFunc {
	return func(c *app.Context) {

		if c.Redis == nil {
			c.JSONErrLog(ecode.ServiceUnavailable("nonce service unavailable"), "nonce service unavailable",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
			)
			c.Abort()
			return
		}

		nonce := c.GetHeader("X-Nonce")
		timestamp := c.GetHeader("X-Timestamp")

		timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			c.JSONErrLog(ecode.Forbidden("timestamp error"), "timestamp parse error",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"timestamp", timestamp,
			)
			c.Abort()
			return
		}

		// 检查时间戳是否过期

		if time.Now().Unix()-timestampInt > 60 {
			c.JSONErrLog(ecode.Forbidden("timestamp expired"), "timestamp expired",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"server_time", time.Now().Unix(),
				"client_time", timestampInt,
			)
			c.Abort()
			return
		}

		// 检查nonce是否已经存在
		noncestr := fmt.Sprintf("Nonce_%s", nonce)

		rval, err := c.Redis.Get(c.Ctx, noncestr)
		if err != nil {
			c.JSONErrLog(ecode.InternalError("nonce error"), "redis get nonce failed",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"key", noncestr,
				"cause", err.Error(),
			)
			c.Abort()
			return
		}

		if rval != "" {
			c.JSONErrLog(ecode.Forbidden("nonce error"), "nonce replay detected",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"key", noncestr,
			)
			c.Abort()
			return
		}

		err = c.Redis.Set(c.Ctx, noncestr, timestamp, 60*time.Second)
		if err != nil {
			c.JSONErrLog(ecode.InternalError("nonce error"), "redis set nonce failed",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"key", noncestr,
				"cause", err.Error(),
			)
			c.Abort()
			return
		}
		c.Next()

	}
}
