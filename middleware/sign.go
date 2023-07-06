package middleware

import (
	"bytes"
	"io/ioutil"
	"sgin/pkg/app"
	"sgin/pkg/utils"
	"sgin/service"
)

// 签名校验中间件
func Signature() app.HandlerFunc {
	return func(c *app.Context) {
		signature := c.GetHeader("X-Signature")
		if signature == "" {
			c.Next()
			return
		}

		appId := c.GetHeader("X-App-Id")
		if appId == "" {
			c.JSONError(403, "X-App-Id is empty")
			c.Abort()
			return
		}

		appInfo, err := service.NewAppService().GetAppByUUID(c, appId)
		if err != nil {
			c.JSONError(403, err.Error())
			c.Abort()
			return
		}

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSONError(403, err.Error())
			return
		}

		// 将 body 内容写回
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		serverSign := utils.SignBody(body, []byte(appInfo.ApiKey))

		if serverSign != signature {
			c.Logger.Error("signature is invalid", "serverSign:", serverSign, "signature:", signature)
			c.JSONError(403, "signature is invalid")
			c.Abort()
			return
		}

		c.Next()
	}
}
