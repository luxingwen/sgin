package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger() HandlerFunc {
	return func(c *Context) {
		// 读取请求体
		// 读取请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			// 处理错误
			c.JSONError(http.StatusBadRequest, "读取请求体失败")
			return
		}

		// 将 body 内容写回
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 限制请求体打印长度，复用 ResponseSize 作为最大打印字节
		reqBody := string(bodyBytes)
		if c.Config != nil && c.Config.LogConfig.ResponseSize > 0 && len(reqBody) > c.Config.LogConfig.ResponseSize {
			reqBody = reqBody[:c.Config.LogConfig.ResponseSize]
		}

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("body", reqBody),
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
		// 基于白名单的跨域控制（支持新旧配置）
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := c.Config.CORS.AllowedOrigins
		if len(allowedOrigins) == 0 {
			allowedOrigins = c.Config.AllowedOrigins
		}

		allowed := false
		if len(allowedOrigins) == 0 {
			// 未配置白名单则回显来源（与旧行为兼容）
			allowed = origin != ""
		} else {
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
		}

		if allowed && origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}

		allowMethods := c.Config.CORS.AllowMethods
		if len(allowMethods) == 0 {
			allowMethods = []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"}
		}
		allowHeaders := c.Config.CORS.AllowHeaders
		if len(allowHeaders) == 0 {
			allowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-Trace-ID", "X-Token", "X-App-Id", "X-Nonce", "X-Timestamp", "X-Signature"}
		}
		exposeHeaders := c.Config.CORS.ExposeHeaders
		if len(exposeHeaders) == 0 {
			exposeHeaders = []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "X-Trace-ID"}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", joinCSV(allowMethods))
		c.Writer.Header().Set("Access-Control-Allow-Headers", joinCSV(allowHeaders))
		if c.Config.CORS.AllowCredentials || len(allowedOrigins) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set("Access-Control-Expose-Headers", joinCSV(exposeHeaders))
		if c.Config.CORS.MaxAge > 0 {
			c.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", c.Config.CORS.MaxAge))
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// 处理norouter
func NoRouterHandler(use ...HandlerFunc) HandlerFunc {
	return func(c *Context) {

		if c.Config.NoRouterFoward == "" {
			c.JSONErrorWithStatus(http.StatusNotFound, http.StatusNotFound, "404 Not Found")
			return
		}

		for _, handler := range use {
			handler(c)
			if c.IsAborted() {
				return
			}
		}

		remote, err := url.Parse(c.Config.NoRouterFoward) //将此替换为你的目标URL
		if err != nil {
			c.Logger.Error(err)
			c.JSONError(http.StatusInternalServerError, "500 Internal Server Error")
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		// 定义我们自己的director
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
		}

		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

func TestAbort() HandlerFunc {
	return func(c *Context) {
		c.AbortWithStatusJSON(http.StatusOK, "test abort")
	}
}

// 设置常见安全响应头
func SecurityHeaders() HandlerFunc {
	return func(c *Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("X-XSS-Protection", "1; mode=block")
		h.Set("Referrer-Policy", "no-referrer")
		// 在非 debug 模式下设置基础 CSP，防止 XSS（可按需放宽/配置化）
		if c.Config == nil || c.Config.LogConfig.Level != "debug" {
			h.Set("Content-Security-Policy", "default-src 'self'")
		}
		c.Next()
	}
}

func joinCSV(items []string) string {
	if len(items) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for i, s := range items {
		if i > 0 {
			buf.WriteByte(',')
			buf.WriteByte(' ')
		}
		buf.WriteString(s)
	}
	return buf.String()
}
