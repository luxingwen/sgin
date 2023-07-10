package middleware

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sgin/pkg/app"
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
				remote, err := url.Parse(c.Config.ForwardAddress) //将此替换为你的目标URL
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
				return
			}
		}

		c.Next()
	}
}
