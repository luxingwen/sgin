package middleware

import (
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sync"

	"golang.org/x/time/rate"
)

// APP 限流中间件
type AppRateLimit struct {
	appLimit map[string]*rate.Limiter
	mu       *sync.RWMutex // 读写锁
	r        rate.Limit    // 令牌桶每秒填充的令牌数
	b        int           // 令牌桶的容量
}

func NewAppRateLimit(r rate.Limit, b int) *AppRateLimit {
	return &AppRateLimit{
		appLimit: make(map[string]*rate.Limiter),
		mu:       &sync.RWMutex{},
		r:        r,
		b:        b,
	}
}

// 获取令牌桶
func (a *AppRateLimit) GetLimit(appId string) *rate.Limiter {
	a.mu.Lock()
	defer a.mu.Unlock()
	l, ok := a.appLimit[appId]
	if !ok {
		l = rate.NewLimiter(a.r, a.b)
		a.appLimit[appId] = l
	}
	return l
}

func (a *AppRateLimit) HandleRateLimit() app.HandlerFunc {
	return func(c *app.Context) {
		// 获取app id
		appId := c.GetString("app_id")

		if appId == "" {
			c.Next()
			return
		}

		// 获取令牌桶
		l := a.GetLimit(appId)

		// 获取令牌
		if !l.Allow() {
			c.JSONErrLog(ecode.TooManyRequests("too many requests"), "too many requests",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"app_id", appId,
			)
			c.Abort()
			return
		}
		c.Next()
	}
}
