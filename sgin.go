package sgin

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luxingwen/sgin/pkg/app"
	"github.com/luxingwen/sgin/pkg/config"
)

// NewApp initializes configuration and returns a new *app.App
// This is a convenience wrapper for external consumers.
func NewApp() *app.App {
	config.InitConfig()
	return app.NewApp()
}

// NewAppFromConfig creates an app using the provided config without calling
// package-level InitConfig. Use this when the host application manages
// configuration itself.
func NewAppFromConfig(cfg *config.Config) *app.App {
	return app.NewAppFromConfig(cfg)
}

// NewAppFromConfigPath loads configuration from the provided path and
// returns a constructed *app.App. Convenience wrapper for host apps.
func NewAppFromConfigPath(path string) *app.App {
	config.InitConfigWithFile(path)
	return app.NewAppFromConfig(config.GetConfig())
}

// RegisterRoutes allows callers to inject their own route registration logic
// into the provided app. Example:
//
//	a := sgin.NewApp()
//	sgin.RegisterRoutes(a, func(a *app.App) { /* custom routes */ })
func RegisterRoutes(a *app.App, fn func(*app.App)) {
	if fn != nil {
		fn(a)
	}
}

// RegisterIntoEngine allows registering sgin plugins/routing callbacks into
// an existing gin Engine from a host application. It constructs a temporary
// App that reuses resources (DB/Logger/Config/Redis) from `src` but sets
// `Router` to the provided `engine`, then calls `fn` so routes attach to the
// host engine.
// Deprecated: RegisterIntoEngine is deprecated. Use RegisterIntoGinEngine when
// embedding into a host that uses `*gin.Engine`. This function remains for
// backward compatibility and will try to forward to the gin-specific helper
// when a `*gin.Engine` is provided; otherwise it falls back to registering
// on the provided `src` App and logs a warning.
func RegisterIntoEngine(src *app.App, engine interface{}, fn func(*app.App)) {
	if src == nil || fn == nil {
		return
	}
	if engine == nil {
		fn(src)
		return
	}
	if ge, ok := engine.(*gin.Engine); ok {
		RegisterIntoGinEngine(src, ge)
		return
	}
	log.Printf("[sgin] RegisterIntoEngine: engine is not *gin.Engine; falling back to registering on src")
	fn(src)
}

// RegisterIntoGinEngine replays plugins registered on `src` into the provided
// `engine` so a host app can reuse routes previously registered with sgin.
// Example:
//
// a := sgin.NewApp()
// a.RegisterPlugin(func(a *app.App) { a.GET("/x", ...) })
// // later, in host with custom engine:
// sgin.RegisterIntoGinEngine(a, engine)
func RegisterIntoGinEngine(src *app.App, engine *gin.Engine) {
	if src == nil || engine == nil {
		return
	}
	temp := &app.App{
		DB:     src.DB,
		Redis:  src.Redis,
		Logger: src.Logger,
		Config: src.Config,
		Router: engine,
	}

	// replay all registered plugins onto the host engine
	for _, p := range src.Plugins {
		if p != nil {
			p(temp)
		}
	}
}

// RegisterPlugin is an alias that forwards plugin registration to App.RegisterPlugin
func RegisterPlugin(a *app.App, fn func(*app.App)) {
	if a == nil || fn == nil {
		return
	}
	a.RegisterPlugin(fn)
}

// Start runs an HTTP server using the provided app's router. If addr is empty,
// it uses the configured ServerPort (prefixed with ':'). This is a convenience
// helper for embedding sgin into other projects.
func Start(a *app.App, addr string) error {
	if a == nil {
		return nil
	}
	if addr == "" {
		if a.Config != nil && a.Config.ServerPort != "" {
			addr = ":" + a.Config.ServerPort
		} else {
			addr = ":8080"
		}
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error(err)
		}
	}()

	// graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}
