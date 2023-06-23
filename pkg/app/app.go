package app

import (
	"sgin/pkg/config"
	"sgin/pkg/db"
	"sgin/pkg/logger"
	"sgin/pkg/redisop"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type App struct {
	DB     *gorm.DB
	Redis  *redisop.RedisClient
	Logger *logger.Logger
	Config *config.Config
	Router *gin.Engine
}

type AppRouterGroup struct {
	*gin.RouterGroup
	App *App
}

func NewApp() *App {
	app := &App{}
	app.Config = config.GetConfig()
	if app.Config.MySQL.Host != "" {
		app.DB = db.GetDB(app.Config.MySQL)
	}

	app.Logger = logger.NewLogger(app.Config.LogConfig)

	if app.Config.MySQL.ShowSQL && app.DB != nil {

		gormLogger := glogger.New(
			app.Logger,
			glogger.Config{
				LogLevel:                  glogger.Info,
				IgnoreRecordNotFoundError: true, // 忽略记录未找到的错误
				Colorful:                  true, // 使用彩色输出
			},
		)

		app.DB.Logger = gormLogger
	}

	app.Router = gin.Default()

	return app
}

func (app *App) Group(relativePath string, handlers ...gin.HandlerFunc) *AppRouterGroup {
	return &AppRouterGroup{
		RouterGroup: app.Router.Group(relativePath, handlers...),
		App:         app,
	}
}

func (app *App) Use(handlers ...HandlerFunc) {
	for _, hf := range handlers {
		app.Router.Use(app.Wrap(hf))
	}
}

func (app *App) GET(relativePath string, hf HandlerFunc) {
	app.Router.GET(relativePath, app.Wrap(hf))
}

func (app *App) POST(relativePath string, hf HandlerFunc) {
	app.Router.POST(relativePath, app.Wrap(hf))
}

func (app *App) PUT(relativePath string, hf HandlerFunc) {
	app.Router.PUT(relativePath, app.Wrap(hf))
}

func (app *App) DELETE(relativePath string, hf HandlerFunc) {
	app.Router.DELETE(relativePath, app.Wrap(hf))
}

func (app *App) PATCH(relativePath string, hf HandlerFunc) {
	app.Router.PATCH(relativePath, app.Wrap(hf))
}

func (rg *AppRouterGroup) GET(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.GET(relativePath, rg.App.Wrap(hf))
}

func (rg *AppRouterGroup) POST(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.POST(relativePath, rg.App.Wrap(hf))
}

func (rg *AppRouterGroup) PUT(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.PUT(relativePath, rg.App.Wrap(hf))
}

func (rg *AppRouterGroup) DELETE(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.DELETE(relativePath, rg.App.Wrap(hf))
}

func (rg *AppRouterGroup) PATCH(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.PATCH(relativePath, rg.App.Wrap(hf))
}

func (rg *AppRouterGroup) Group(relativePath string, handlers ...gin.HandlerFunc) *AppRouterGroup {
	return &AppRouterGroup{
		RouterGroup: rg.RouterGroup.Group(relativePath, handlers...),
		App:         rg.App,
	}
}

func (rg *AppRouterGroup) Use(handlers ...HandlerFunc) {
	for _, hf := range handlers {
		rg.RouterGroup.Use(rg.App.Wrap(hf))
	}
}

// Similarly, define other HTTP method handlers like POST, PUT, DELETE...
