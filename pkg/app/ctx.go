package app

import (
	"context"

	"github.com/luxingwen/sgin/pkg/config"
	"github.com/luxingwen/sgin/pkg/logger"
	"github.com/luxingwen/sgin/pkg/redisop"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Context struct {
	*gin.Context
	DB      *gorm.DB
	Redis   *redisop.RedisClient
	Logger  *logger.Logger
	Config  *config.Config
	TraceID string
	Ctx     context.Context
}

type HandlerFunc func(*Context)

// AppContext 定义了应用层可用的最小上下文接口。
// 保持向后兼容：现有的 *Context 会实现该接口。
type AppContext interface {
	GetDB() *gorm.DB
	GetRedis() *redisop.RedisClient
	GetLogger() *logger.Logger
	GetConfig() *config.Config
	GetTraceID() string
	GetCtx() context.Context
	GinContext() *gin.Context // 非 HTTP 场景返回 nil
}

// 新的以接口为参数的 handler，方便宿主以接口编程逐步迁移
type HandlerFuncIface func(AppContext)

// WrapIface 把以 AppContext 为参数的处理函数包装为 gin.HandlerFunc
func (app *App) WrapIface(hf HandlerFuncIface) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("X-Trace-ID")

		if traceID == "" {
			traceID = uuid.New().String()
		}

		c.Writer.Header().Set("X-Trace-ID", traceID)

		cc := &Context{
			Context: c,
			DB:      app.DB,
			Redis:   app.Redis,
			Logger: app.Logger.With(
				zap.String("traceID", traceID),
			),
			Config:  app.Config,
			TraceID: traceID,
			Ctx:     c.Request.Context(),
		}
		hf(cc)
	}
}

func (app *App) Wrap(hf HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("X-Trace-ID")

		if traceID == "" {
			traceID = uuid.New().String()
		}

		// ensure trace id is visible to clients
		c.Writer.Header().Set("X-Trace-ID", traceID)

		cc := &Context{
			Context: c,
			DB:      app.DB,
			Redis:   app.Redis,
			Logger: app.Logger.With(
				zap.String("traceID", traceID),
			),
			Config:  app.Config,
			TraceID: traceID,
			Ctx:     c.Request.Context(),
		}
		hf(cc)
	}
}

// AppContext 接口的实现 — 让当前的 Context 满足 AppContext
func (c *Context) GetDB() *gorm.DB                { return c.DB }
func (c *Context) GetRedis() *redisop.RedisClient { return c.Redis }
func (c *Context) GetLogger() *logger.Logger      { return c.Logger }
func (c *Context) GetConfig() *config.Config      { return c.Config }
func (c *Context) GetTraceID() string             { return c.TraceID }
func (c *Context) GetCtx() context.Context        { return c.Ctx }
func (c *Context) GinContext() *gin.Context       { return c.Context }
