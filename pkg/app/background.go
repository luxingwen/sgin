package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luxingwen/sgin/pkg/config"
	"github.com/luxingwen/sgin/pkg/logger"
	"github.com/luxingwen/sgin/pkg/redisop"
	"gorm.io/gorm"
)

// BackgroundContext 是一个不依赖 gin 的 AppContext 实现，适用于命令行、cron、后台任务等场景。
type BackgroundContext struct {
	DB      *gorm.DB
	Redis   *redisop.RedisClient
	Logger  *logger.Logger
	Config  *config.Config
	TraceID string
	Ctx     context.Context
}

// NewBackgroundContextFromApp 从现有 *App 构建 BackgroundContext（复用 App 的 DB/Logger/Redis/Config）。
func NewBackgroundContextFromApp(a *App) *BackgroundContext {
	if a == nil {
		return &BackgroundContext{TraceID: uuid.New().String(), Ctx: context.Background()}
	}
	trace := uuid.New().String()
	return &BackgroundContext{
		DB:      a.DB,
		Redis:   a.Redis,
		Logger:  a.Logger.With(),
		Config:  a.Config,
		TraceID: trace,
		Ctx:     context.Background(),
	}
}

// NewBackgroundContext 构造一个 BackgroundContext，允许直接传入依赖。
func NewBackgroundContext(db *gorm.DB, redis *redisop.RedisClient, logger *logger.Logger, cfg *config.Config) *BackgroundContext {
	return &BackgroundContext{
		DB:      db,
		Redis:   redis,
		Logger:  logger,
		Config:  cfg,
		TraceID: uuid.New().String(),
		Ctx:     context.Background(),
	}
}

// 实现 AppContext 接口
func (b *BackgroundContext) GetDB() *gorm.DB                { return b.DB }
func (b *BackgroundContext) GetRedis() *redisop.RedisClient { return b.Redis }
func (b *BackgroundContext) GetLogger() *logger.Logger      { return b.Logger }
func (b *BackgroundContext) GetConfig() *config.Config      { return b.Config }
func (b *BackgroundContext) GetTraceID() string             { return b.TraceID }
func (b *BackgroundContext) GetCtx() context.Context        { return b.Ctx }
func (b *BackgroundContext) GinContext() *gin.Context       { return nil }
