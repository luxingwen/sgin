package main

import (
	"sgin/pkg/app"
	"sgin/pkg/config"
)

func main() {
	config.InitConfig()
	serverApp := app.NewApp()
	serverApp.Use(app.RecoveryWithWriter(serverApp.Logger))

	serverApp.GET("/ping", func(ctx *app.Context) {
		panic("test panic")
		ctx.JSONSuccess("pong")
	})

	serverApp.Router.Run(":8080")
}
