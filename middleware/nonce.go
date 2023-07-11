package middleware

import (
	"fmt"
	"net/http"
	"sgin/pkg/app"
	"strconv"
	"time"
)

// 防重放攻击中间件
func NonceHandler() app.HandlerFunc {
	return func(c *app.Context) {

		nonce := c.GetHeader("X-Nonce")
		timestamp := c.GetHeader("X-Timestamp")

		timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			c.JSONError(http.StatusForbidden, "timestamp error")
			c.Abort()
			return
		}

		// 检查时间戳是否过期

		if time.Now().Unix()-timestampInt > 60 {
			c.Logger.Error("timestamp expired", "server time", time.Now().Unix(), "client time", timestampInt)
			c.JSONError(http.StatusForbidden, "timestamp expired")
			c.Abort()
			return
		}

		// 检查nonce是否已经存在
		noncestr := fmt.Sprintf("Nonce_%s", nonce)

		rval, err := c.Redis.Get(c, noncestr)
		if err != nil {
			c.JSONError(http.StatusForbidden, "nonce error")
			c.Abort()
			return
		}

		if rval != "" {
			c.JSONError(http.StatusForbidden, "nonce error")
			c.Abort()
			return
		}

		err = c.Redis.Set(c, noncestr, timestamp, 60*time.Second)
		if err != nil {
			c.JSONError(http.StatusForbidden, "nonce error")
			c.Abort()
			return
		}
		c.Next()

	}
}
