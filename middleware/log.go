package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 日志J=记录中间件
func LogMiddleware() app.HandlerFunc {
	return func(c *app.Context) {

		// 获取请求信息

		header := c.Request.Header
		headerByte, _ := json.Marshal(header)
		ip := c.ClientIP()

		path := c.Request.URL.Path
		method := c.Request.Method

		// 读取请求体
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			// 处理错误
			c.JSONError(http.StatusBadRequest, "读取请求体失败")
			return
		}

		// 将 body 内容写回
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		userId := c.GetString("user_id")
		appid := c.GetHeader("X-App-Id")
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
		}

		// 将日志信息写入数据库
		err = service.NewLogService().CreateLog(c, &logInfo)
		if err != nil {
			c.JSONError(http.StatusBadRequest, err.Error())
			return
		}

		c.Next()

		// 获取响应信息
		status := c.Writer.Status()
		logInfo.Status = status

		// 更新日志
		err = service.NewLogService().UpdateLog(c, &logInfo)
		if err != nil {
			c.Logger.Error("update log error", zap.Error(err))
		}

	}
}
