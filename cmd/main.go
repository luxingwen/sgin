package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sgin/pkg/app"
	"sgin/pkg/config"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()
	serverApp := app.NewApp()
	serverApp.Use(app.RecoveryWithWriter(serverApp.Logger))

	serverApp.GET("/ping", func(ctx *app.Context) {
		panic("test panic")
		ctx.JSONSuccess("pong")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: serverApp.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverApp.Logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	serverApp.Logger.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		serverApp.Logger.Fatal("Server forced to shutdown: ", err)
	}

	serverApp.Logger.Info("Server exiting")
}
