package middleware

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
)

// 根据前缀进行转发
func ForwardByPrefix(use ...app.HandlerFunc) app.HandlerFunc {
	return func(c *app.Context) {

		upath := c.Request.URL.Path
		for _, prefix := range c.Config.ForwardPrefix {
			if len(upath) >= len(prefix) && upath[:len(prefix)] == prefix {

				for _, handler := range use {
					handler(c)
					if c.IsAborted() {
						return
					}
				}
				// 转发
				remote, err := url.Parse(c.Config.ForwardAddress)
				if err != nil {
					c.JSONErrLog(ecode.InternalError("forward address invalid"), "parse forward address failed",
						"trace_id", c.TraceID,
						"path", c.FullPath(),
						"method", c.Request.Method,
						"client_ip", c.ClientIP(),
						"forward_address", c.Config.ForwardAddress,
						"cause", err.Error(),
					)
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
				return
			}
		}

		c.Next()
	}
}
