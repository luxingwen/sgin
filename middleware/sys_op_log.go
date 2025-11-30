package middleware

import (
	"bytes"
	"io"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
	"strings"
	"time"
)

// 定义要过滤的路径
var opFilterPath = []string{
	"/api/v1/login",
	"/api/v1/verification_code/create",
	"/api/v1/upload/file/*",
	"/ping",
	"/public/",
	"/swagger/",
	"/api/v1/sysoplog/list",
	"/api/v1/user/myinfo",
}

// SysOpLogMiddleware 记录操作日志的中间件
func SysOpLogMiddleware(logservice *service.SysOpLogService) app.HandlerFunc {
	return func(c *app.Context) {
		startTime := time.Now()

		ip := c.ClientIP()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 过滤不需要记录的路径
		for _, filterPath := range opFilterPath {
			// 获取最后一个字符是否为*，如果是则表示前缀匹配
			if filterPath[len(filterPath)-1] == '*' {
				if len(path) >= len(filterPath)-1 && path[:len(filterPath)-1] == filterPath[:len(filterPath)-1] {
					c.Next()
					return
				}
			} else {
				if path == filterPath {
					c.Next()
					return
				}
			}
		}

		bodyBytes := []byte{}
		var err error

		if method == "POST" {
			// 获取content-type
			contentType := c.GetHeader("Content-Type")
			if strings.HasPrefix(contentType, "application/json") {
				// 读取请求体
				bodyBytes, err = io.ReadAll(c.Request.Body)
				if err != nil {
					c.JSONErrLog(ecode.BadRequest("读取请求体失败"), "read request body failed",
						"trace_id", c.TraceID,
						"path", c.FullPath(),
						"method", c.Request.Method,
						"client_ip", c.ClientIP(),
						"cause", err.Error(),
					)
					return
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		userId := c.GetString("user_id")

		logInfo := model.SysOpLog{
			UserUuid:  userId,
			Path:      path,
			Method:    method,
			Ip:        ip,
			Status:    0, // 初始状态，后面更新
			Message:   "",
			Params:    string(bodyBytes),
			Duration:  0, // 初始时为0，后面更新
			CreatedAt: startTime.Format("2006-01-02 15:04:05"),
			RequestId: c.TraceID,
		}

		c.Next()

		// 获取响应信息
		status := c.Writer.Status()
		msg := c.GetString("message")
		code := c.GetInt("code")
		logInfo.Status = status
		logInfo.Message = msg
		logInfo.Code = code
		logInfo.Duration = time.Since(startTime).Milliseconds()
		// logInfo.Response = string(c.Writer.Body.Bytes()) // 获取响应体

		// 将日志信息写入数据库
		err = logservice.CreateSysOpLog(c, &logInfo)
		if err != nil {
			c.Logger.Errorw("create sys op log failed",
				"trace_id", c.TraceID,
				"path", c.FullPath(),
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"cause", err.Error(),
			)
			return
		}
	}
}
