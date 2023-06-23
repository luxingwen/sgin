package app

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger() HandlerFunc {
	return func(c *Context) {
		// 读取请求体
		// 读取请求体
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			// 处理错误
			c.JSONError(http.StatusBadRequest, "读取请求体失败")
			return
		}

		// 将 body 内容写回
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("body", string(bodyBytes)),
		}
		c.Logger.WithOptions(zap.Fields(fields...)).Info("Request")

		c.Next()
	}
}

type ResponseRecorder struct {
	gin.ResponseWriter
	Body bytes.Buffer
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *ResponseRecorder) WriteString(s string) (int, error) {
	r.Body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}

func ResponseLogger() HandlerFunc {
	return func(c *Context) {

		recorder := &ResponseRecorder{ResponseWriter: c.Writer}
		c.Writer = recorder

		c.Next()

		c.Logger.With(zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path), zap.String("ip", c.ClientIP()), zap.Int("status", c.Writer.Status())).Info("Response")

		// 读取响应体

		if c.Config.LogConfig.ResponseSize > 0 {
			body := recorder.Body.Bytes()
			if len(body) > c.Config.LogConfig.ResponseSize {
				body = body[:c.Config.LogConfig.ResponseSize]
			}
			c.Logger.Infof("Response body: %s", string(body))
		}

	}
}

// 处理跨域请求
func Cors() HandlerFunc {
	return func(c *Context) {
		// 设置跨域请求头
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Trace-ID")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
