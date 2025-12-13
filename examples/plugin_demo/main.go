package main

import (
	"fmt"
	"log"

	"github.com/luxingwen/sgin/examples/plugin_demo/plugin"
	"github.com/luxingwen/sgin/pkg/app"
	"github.com/luxingwen/sgin/pkg/config"
)

func main() {
	// initialize config from example config file
	config.InitConfigWithFile("examples/plugin_demo/config.yaml")

	// show that plugin config was decoded by RegisterExtension
	if plugin.Conf != nil {
		fmt.Printf("plugin config: feature=%v, limit=%d\n", plugin.Conf.FeatureFlag, plugin.Conf.Limit)
	} else {
		fmt.Println("plugin config not found or failed to parse")
	}

	// create App from the loaded config
	a := app.NewAppFromConfig(config.GetConfig())

	// let plugin register routes / middleware on the app
	plugin.Setup(a)

	// add a simple ping endpoint
	a.GET("/ping", func(c *app.Context) {
		c.JSONSuccess("pong from example")
	})

	// run server for demonstration
	log.Println("starting example server at :8080 (use Ctrl+C to stop)")
	if err := a.Router.Run(":8080"); err != nil {
		log.Fatalf("server exit: %v", err)
	}
}
