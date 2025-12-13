package app

import (
	"github.com/luxingwen/sgin/pkg/config"
	"github.com/luxingwen/sgin/pkg/db"
	"github.com/luxingwen/sgin/pkg/logger"
	"github.com/luxingwen/sgin/pkg/redisop"

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
	// Plugins stores registered plugin callbacks so they can be replayed
	// into a different router (useful for embedding into a host engine).
	Plugins []func(*App)
}

// RegisterPlugin 允许宿主或外部模块以回调方式注册路由/中间件等
func (a *App) RegisterPlugin(fn func(*App)) {
	if fn == nil {
		return
	}
	// store plugin for later replay and invoke immediately for existing behavior
	a.Plugins = append(a.Plugins, fn)
	fn(a)
}

// StorePlugin stores a plugin callback without invoking it. This is useful
// for embedding scenarios where the host wants to replay registered plugins
// into its own router instead of immediately applying them to `a.Router`.
func (a *App) StorePlugin(fn func(*App)) {
	if fn == nil {
		return
	}
	a.Plugins = append(a.Plugins, fn)
}

type AppRouterGroup struct {
	*gin.RouterGroup
	App *App
}

func NewApp() *App {
	return NewAppFromConfig(config.GetConfig())
}

// NewAppFromConfig creates an App using the provided Config. This is useful
// when the host application already has its own configuration and does not
// want sgin to call InitConfig/read files itself.
func NewAppFromConfig(cfg *config.Config) *App {
	a := &App{}
	a.Config = cfg
	if a.Config == nil {
		// caller likely forgot to init config; fall back to default behavior
		a.Config = config.GetConfig()
	}

	if a.Config.MySQL.Host != "" {
		a.DB = db.GetDB(a.Config.MySQL)
	}

	a.Logger = logger.NewLogger(a.Config.LogConfig)

	if a.Config.MySQL.ShowSQL && a.DB != nil {
		gormLogger := glogger.New(
			a.Logger,
			glogger.Config{
				LogLevel:                  glogger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
		a.DB.Logger = gormLogger
	}

	if a.Config.RedisConfig.Address != "" {
		a.Redis = redisop.NewRedisClient(a.Config.RedisConfig.Address, a.Config.RedisConfig.Password, a.Config.RedisConfig.Database)
	}

	if a.Config.LogConfig.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	a.Router = gin.New()

	return a
}

// NewAppFromConfigPath loads configuration from the given file path and
// returns a newly constructed *App. This is useful when a host wants to
// provide a specific config file path without manipulating env vars.
func NewAppFromConfigPath(path string) *App {
	config.InitConfigWithFile(path)
	return NewAppFromConfig(config.GetConfig())
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

func (app *App) NoRoute(hf HandlerFunc) {
	app.Router.NoRoute(app.Wrap(hf))
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
