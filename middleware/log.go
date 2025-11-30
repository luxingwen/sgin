package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	// 限制请求与响应体记录大小，避免占用过多内存/存储
	maxReqBodyLogBytes  = 1 << 20 // 1MB
	maxRespBodyLogBytes = 1 << 20 // 1MB
)

type bodyWriter struct {
	gin.ResponseWriter
	buf bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	// 限制记录的响应体大小
	if w.buf.Len() < maxRespBodyLogBytes {
		// 只追加未超过上限的部分
		remaining := maxRespBodyLogBytes - w.buf.Len()
		if len(b) > remaining {
			w.buf.Write(b[:remaining])
		} else {
			w.buf.Write(b)
		}
	}
	return w.ResponseWriter.Write(b)
}

func (w *bodyWriter) WriteString(s string) (int, error) {
	if w.buf.Len() < maxRespBodyLogBytes {
		remaining := maxRespBodyLogBytes - w.buf.Len()
		if len(s) > remaining {
			w.buf.WriteString(s[:remaining])
		} else {
			w.buf.WriteString(s)
		}
	}
	return w.ResponseWriter.WriteString(s)
}

// 日志记录中间件（最佳努力，不影响业务返回）
func LogMiddleware() app.HandlerFunc {
	return func(c *app.Context) {
		// 获取请求信息
		// 复制并脱敏敏感头
		header := c.Request.Header.Clone()
		if header.Get("Authorization") != "" {
			header.Set("Authorization", "*****")
		}
		if header.Get("X-Token") != "" {
			header.Set("X-Token", "*****")
		}
		if header.Get("X-Signature") != "" {
			header.Set("X-Signature", "*****")
		}
		if header.Get("Cookie") != "" {
			header.Set("Cookie", "*****")
		}
		headerByte, _ := json.Marshal(header)
		ip := c.ClientIP()

		path := c.Request.URL.Path
		method := c.Request.Method

		// 读取请求体（带大小限制，避免耗尽内存；二进制或表单略过）
		var bodyBytes []byte
		var err error
		ct := c.GetHeader("Content-Type")
		if ct != "" && (contains(ct, "multipart/form-data") || contains(ct, "application/octet-stream")) {
			bodyBytes = []byte("[binary body omitted]")
		} else {
			// 限流读取
			limited := io.LimitReader(c.Request.Body, maxReqBodyLogBytes)
			bodyBytes, err = io.ReadAll(limited)
			if err != nil {
				// 读取失败不影响主流程
				c.Logger.Warnw("read request body failed", "error", err.Error(), "trace_id", c.TraceID, "path", path, "method", method)
			}
		}
		// 将 body 内容写回
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		userId := c.GetString("user_id")
		appid := c.GetString("app_id")
		if appid == "" {
			appid = c.GetHeader("X-App-Id")
		}

		logInfo := model.Log{
			Header:   string(headerByte),
			Ip:       ip,
			Path:     path,
			Method:   method,
			ReqBody:  string(bodyBytes),
			AppId:    appid,
			UserUUID: userId,
			TraceID:  c.TraceID,
			UUID:     uuid.New().String(),
			Action:   c.FullPath(),
		}

		// 将日志信息写入数据库
		if err = service.NewLogService().CreateLog(c, &logInfo); err != nil {
			// 不中断业务，仅记录错误
			c.Logger.Errorw("create request log failed",
				"error", err.Error(),
				"trace_id", c.TraceID,
				"path", path,
				"method", method,
				"client_ip", ip,
				"app_id", appid,
				"user_id", userId,
			)
		}

		// 包装 ResponseWriter 以捕获响应体
		bw := &bodyWriter{ResponseWriter: c.Writer}
		c.Writer = bw

		c.Next()

		// 获取响应信息
		status := c.Writer.Status()
		respBody := bw.buf.String()
		// 从 context 中获取 message（由统一响应设置）
		if msgVal, exists := c.Get("message"); exists {
			if s, ok := msgVal.(string); ok {
				logInfo.Message = s
			}
		}
		logInfo.Status = status
		logInfo.RespBody = respBody

		// 更新日志
		if err = service.NewLogService().UpdateLog(c, &logInfo); err != nil {
			c.Logger.Error("update log error", zap.Error(err))
		}

	}
}

// contains 判断子串是否在字符串中（大小写不敏感场景此处保持原值判断）
func contains(s, sub string) bool {
	return bytes.Contains([]byte(s), []byte(sub))
}
