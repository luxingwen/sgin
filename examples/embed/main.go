package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/luxingwen/sgin"
	"github.com/luxingwen/sgin/pkg/app"
	_ "github.com/luxingwen/sgin/pkg/config"
	"github.com/luxingwen/sgin/routers"
)

func main() {
	// support a configurable config file path via -config
	cfgPath := flag.String("config", "", "path to config file (overrides CONFIG_FILE env)")
	// addr := flag.String("addr", ":8080", "address for host engine to listen on")
	flag.Parse()

	if *cfgPath != "" {
		// set env var expected by pkg/config.InitConfig
		os.Setenv("CONFIG_FILE", *cfgPath)
		fmt.Printf("using config file: %s\n", *cfgPath)
	}

	// create app (calls InitConfig internally)
	a := sgin.NewApp()
	// store the full router set as plugin callbacks (do NOT register on a.Router)
	routers.InitRouterStored(a)

	// register plugin / custom routes on sgin App (this will also store
	// the plugin so it can be replayed into a host engine)
	a.RegisterPlugin(func(a *app.App) {
		a.GET("/hello", func(ctx *app.Context) {
			ctx.JSONSuccess("hello from embedded sgin")
		})
	})

	// Host application creates its own gin engine and can replay the
	// previously-registered sgin plugins into the host engine so routes
	// are served by the host instead of sgin's internal engine.
	//engine := gin.New()
	// (attach host middleware, logging, etc.)

	// Replay sgin's registered plugins into host engine
	//sgin.RegisterIntoGinEngine(a, engine)

	// Build a host App that reuses sgin resources but serves using the
	// host engine, then start it using sgin's Start helper which provides
	// graceful shutdown behavior.
	// hostApp := &app.App{
	// 	DB:     a.DB,
	// 	Redis:  a.Redis,
	// 	Logger: a.Logger,
	// 	Config: a.Config,
	// 	Router: engine,
	// }

	_ = sgin.Start(a, "")
}
