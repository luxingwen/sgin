package app

import (
	"sgin/pkg/config"
	"sgin/pkg/logger"
	"sgin/pkg/redisop"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"net/http"

	"go.uber.org/zap"
)

type WSContext struct {
	Conn    *websocket.Conn
	DB      *gorm.DB
	Redis   *redisop.RedisClient
	Logger  *logger.Logger
	Config  *config.Config
	TraceID string
}

type WsHandlerFunc func(*WSContext)

func (app *App) WrapWS(hf WsHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			app.Logger.Error(err)
			return
		}

		traceID := r.Header.Get("X-Trace-ID")

		if traceID == "" {
			traceID = uuid.New().String()
		}

		cc := &WSContext{
			Conn:  conn,
			DB:    app.DB,
			Redis: app.Redis,
			Logger: app.Logger.With(
				zap.String("traceID", traceID),
			),
			Config:  app.Config,
			TraceID: traceID,
		}
		hf(cc)
	}
}
