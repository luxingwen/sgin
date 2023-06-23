package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sgin/pkg/config"
	"sgin/pkg/logger"
	"sgin/pkg/redisop"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Context struct {
	*gin.Context
	DB      *gorm.DB
	Redis   *redisop.RedisClient
	Logger  *logger.Logger
	Config  *config.Config
	TraceID string
	Ctx     context.Context
}

type HandlerFunc func(*Context)

func (app *App) Wrap(hf HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("X-Trace-ID")

		if traceID == "" {
			traceID = uuid.New().String()
		}

		cc := &Context{
			Context: c,
			DB:      app.DB,
			Redis:   app.Redis,
			Logger: app.Logger.With(
				zap.String("traceID", traceID),
			),
			Config:  app.Config,
			TraceID: traceID,
		}
		hf(cc)
	}
}

type Response struct {
	TraceID string      `json:"trace_id"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (ctx *Context) JSONSuccess(data interface{}) {
	response := Response{
		TraceID: ctx.TraceID,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	}

	ctx.JSON(http.StatusOK, response)
}

func (ctx *Context) JSONError(code int, message string) {
	response := Response{
		TraceID: ctx.TraceID,
		Code:    code,
		Message: message,
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, response)
}

func (ctx *Context) ReturnWithStream_ObjTxt(data interface{}) {
	response := Response{
		TraceID: ctx.TraceID,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	}
	dataByte, _ := json.Marshal(response)
	respData := string(dataByte)
	var fullData []byte
	fullData = append(fullData, respData...)
	fullData = append(fullData, []byte("\n\n")...)
	resp := bytes.NewBuffer(fullData)
	fmt.Fprintf(ctx.Writer, "%s", resp)
	ctx.Writer.(http.Flusher).Flush()
	return
}
