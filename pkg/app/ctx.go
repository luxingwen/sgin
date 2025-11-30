package app

import (
	"context"
	"sgin/pkg/config"
	"sgin/pkg/logger"
	"sgin/pkg/redisop"

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
