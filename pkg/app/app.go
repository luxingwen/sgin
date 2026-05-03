package app

import (
	"github.com/luxingwen/sgin/pkg/config"
	"github.com/luxingwen/sgin/pkg/db"
	"github.com/luxingwen/sgin/pkg/logger"
	"github.com/luxingwen/sgin/pkg/redisop"

	"net/http"
	"reflect"

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
	// BasePath is an optional prefix automatically added to all routes registered via App methods
	BasePath string
	// Plugins stores registered plugin callbacks so they can be replayed
	// into a different router (useful for embedding into a host engine).
	Plugins []func(*App)
	// Extras holds decoded custom configuration structures keyed by caller-provided name
	Extras map[string]interface{}
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

// GetExtra returns the extra object registered under name, if any.
func (a *App) GetExtra(name string) (interface{}, bool) {
	if a == nil || a.Extras == nil {
		return nil, false
	}
	v, ok := a.Extras[name]
	return v, ok
}

// GetExtraAs attempts to return the extra stored under name as type T.
// It supports values stored as either T or *T; when stored as *T and T is
// a non-pointer type, the pointed value will be returned.
func GetExtraAs[T any](a *App, name string) (T, bool) {
	var zero T
	if a == nil {
		return zero, false
	}
	v, ok := a.GetExtra(name)
	if !ok {
		return zero, false
	}
	// direct assertion
	if val, ok := v.(T); ok {
		return val, true
	}
	// if v is a pointer to T, dereference and try
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && rv.IsValid() {
		ev := rv.Elem()
		if ev.IsValid() {
			if iface := ev.Interface(); iface != nil {
				if val, ok := iface.(T); ok {
					return val, true
				}
			}
		}
	}
	return zero, false
}

type AppRouterGroup struct {
	*gin.RouterGroup
	App *App
}

func NewApp() *App {
	return NewAppFromConfig(config.GetConfig())
}

// AppOption allows callers to customize App during initialization.
type AppOption func(*App)

// WithExtra decodes configuration key into the provided out pointer and stores it
// in App.Extras under the provided name. Example:
//
//	app.NewAppWithOptions(app.WithExtra("my_extra", "my_extra", &MyExtra{}))
func WithExtra(name string, key string, out interface{}) AppOption {
	return func(a *App) {
		// best-effort: try to unmarshal the key into out; ignore error here
		_ = config.UnmarshalKey(key, out)
		if a.Extras == nil {
			a.Extras = make(map[string]interface{})
		}
		a.Extras[name] = out
	}
}

// NewAppWithOptions creates a new App and applies the provided options.
func NewAppWithOptions(opts ...AppOption) *App {
	a := NewApp()
	for _, o := range opts {
		if o != nil {
			o(a)
		}
	}
	return a
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

	if a.Config.DBType == "postgres" && a.Config.Postgres.Host != "" {
		a.DB = db.GetDB("postgres", a.Config.Postgres)
	} else if a.Config.MySQL.Host != "" {
		a.DB = db.GetDB("mysql", a.Config.MySQL)
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
		RouterGroup: app.Router.Group(app.addBase(relativePath), handlers...),
		App:         app,
	}
}

func (app *App) Use(handlers ...HandlerFunc) {
	for _, hf := range handlers {
		app.Router.Use(app.Wrap(hf))
	}
}

func (app *App) GET(relativePath string, hf HandlerFunc) {
	app.Router.GET(app.addBase(relativePath), app.Wrap(hf))
}

func (app *App) POST(relativePath string, hf HandlerFunc) {
	app.Router.POST(app.addBase(relativePath), app.Wrap(hf))
}

func (app *App) PUT(relativePath string, hf HandlerFunc) {
	app.Router.PUT(app.addBase(relativePath), app.Wrap(hf))
}

func (app *App) DELETE(relativePath string, hf HandlerFunc) {
	app.Router.DELETE(app.addBase(relativePath), app.Wrap(hf))
}

func (app *App) PATCH(relativePath string, hf HandlerFunc) {
	app.Router.PATCH(app.addBase(relativePath), app.Wrap(hf))
}

func (app *App) NoRoute(hf HandlerFunc) {
	app.Router.NoRoute(app.Wrap(hf))
}

// Static 在指定相对路径处提供静态文件
func (app *App) Static(relativePath string, root string) {
	app.Router.Static(app.addBase(relativePath), root)
}

// StaticFS 在指定相对路径处提供来自http.FileSystem的静态文件
func (app *App) StaticFS(relativePath string, fs http.FileSystem) {
	app.Router.StaticFS(app.addBase(relativePath), fs)
}

// StaticFile 在指定相对路径处提供单个静态文件
func (app *App) StaticFile(relativePath string, filePath string) {
	app.Router.StaticFile(app.addBase(relativePath), filePath)
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

// Static 在指定相对路径处提供静态文件
func (rg *AppRouterGroup) Static(relativePath string, root string) {
	rg.RouterGroup.Static(relativePath, root)
}

// StaticFS 在指定相对路径处提供来自http.FileSystem的静态文件
func (rg *AppRouterGroup) StaticFS(relativePath string, fs http.FileSystem) {
	rg.RouterGroup.StaticFS(relativePath, fs)
}

// StaticFile 在指定相对路径处提供单个静态文件
func (rg *AppRouterGroup) StaticFile(relativePath string, filePath string) {
	rg.RouterGroup.StaticFile(relativePath, filePath)
}

// SetBasePath sets a prefix applied to all routes registered on App (and App.Group).
// Prefix should start with "/"; trailing slashes are trimmed to avoid double slashes.
func (app *App) SetBasePath(prefix string) {
	if prefix == "" || prefix == "/" {
		app.BasePath = ""
		return
	}
	// normalize to "/something" without trailing slash
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}
	// trim trailing slash except root
	for len(prefix) > 1 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	app.BasePath = prefix
}

// addBase prefixes path with BasePath when configured.
func (app *App) addBase(relativePath string) string {
	if app == nil || app.BasePath == "" {
		return relativePath
	}
	// ensure relativePath starts with "/"
	if relativePath == "" || relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	return app.BasePath + relativePath
}

// Similarly, define other HTTP method handlers like POST, PUT, DELETE...
