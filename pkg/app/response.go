package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sgin/pkg/ecode"
)

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
	ctx.Set("code", http.StatusOK)
	ctx.Set("message", "Success")
	ctx.JSON(http.StatusOK, response)
}

func (ctx *Context) JSONError(code int, message string) {
	response := Response{
		TraceID: ctx.TraceID,
		Code:    code,
		Message: message,
		Data:    nil,
	}
	ctx.Set("code", code)
	ctx.Set("message", message)

	ctx.JSON(http.StatusOK, response)
}

// JSONErrorWithStatus 返回携带真实 HTTP 状态码的错误响应
func (ctx *Context) JSONErrorWithStatus(httpStatus int, code int, message string) {
	response := Response{
		TraceID: ctx.TraceID,
		Code:    code,
		Message: message,
		Data:    nil,
	}
	ctx.Set("code", code)
	ctx.Set("message", message)
	ctx.JSON(httpStatus, response)
}

// Convenience helpers with standard HTTP semantics
func (ctx *Context) JSONBadRequest(message string) {
	if message == "" {
		message = "Bad Request"
	}
	ctx.JSONErrorWithStatus(http.StatusBadRequest, http.StatusBadRequest, message)
}

func (ctx *Context) JSONUnauthorized(message string) {
	if message == "" {
		message = "Unauthorized"
	}
	ctx.JSONErrorWithStatus(http.StatusUnauthorized, http.StatusUnauthorized, message)
}

func (ctx *Context) JSONForbidden(message string) {
	if message == "" {
		message = "Forbidden"
	}
	ctx.JSONErrorWithStatus(http.StatusForbidden, http.StatusForbidden, message)
}

func (ctx *Context) JSONNotFound(message string) {
	if message == "" {
		message = "Not Found"
	}
	ctx.JSONErrorWithStatus(http.StatusNotFound, http.StatusNotFound, message)
}

func (ctx *Context) JSONConflict(message string) {
	if message == "" {
		message = "Conflict"
	}
	ctx.JSONErrorWithStatus(http.StatusConflict, http.StatusConflict, message)
}

func (ctx *Context) JSONInternalError(message string) {
	if message == "" {
		message = "Internal Server Error"
	}
	ctx.JSONErrorWithStatus(http.StatusInternalServerError, http.StatusInternalServerError, message)
}

// JSONErr 根据 ecode.APIError 或普通 error 统一响应
func (ctx *Context) JSONErr(err error) {
	if err == nil {
		ctx.JSONSuccess(nil)
		return
	}
	if ae, ok := err.(*ecode.APIError); ok {
		// 为兼容现有客户端，默认仍以 200 HTTP 状态返回，错误码置于 body.code
		ctx.JSONError(ae.Code, ae.Message)
		return
	}
	ctx.JSONError(http.StatusInternalServerError, "Internal Server Error")
}

// JSONErrLog 记录错误日志并返回统一 JSON 错误
// 约定：4xx 使用 Warn，5xx 使用 Error；仍以 HTTP 200 返回（兼容现有客户端）
func (ctx *Context) JSONErrLog(err error, msg string, keysAndValues ...interface{}) {
	if err == nil {
		ctx.JSONSuccess(nil)
		return
	}
	// 基础上下文字段
	base := []interface{}{
		"trace_id", ctx.TraceID,
		"path", ctx.FullPath(),
		"method", ctx.Request.Method,
		"client_ip", ctx.ClientIP(),
	}
	if uid := ctx.GetString("user_id"); uid != "" {
		base = append(base, "user_id", uid)
	}
	if aid := ctx.GetString("app_id"); aid != "" {
		base = append(base, "app_id", aid)
	}
	// 合并用户传入字段
	fields := append(base, keysAndValues...)

	if ae, ok := err.(*ecode.APIError); ok {
		fields = append(fields, "code", ae.Code)
		// 4xx -> Warn，5xx -> Error
		if ae.HTTPStatus >= 500 {
			ctx.Logger.Errorw(msg, append(fields, "error", ae.Message)...)
		} else {
			ctx.Logger.Warnw(msg, append(fields, "error", ae.Message)...)
		}
		ctx.JSONError(ae.Code, ae.Message)
		return
	}
	ctx.Logger.Errorw(msg, append(fields, "error", err.Error())...)
	ctx.JSONError(http.StatusInternalServerError, "Internal Server Error")
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
}

// Response data with list object
func (ctx *Context) ResList(v interface{}) {
	ctx.JSONSuccess(ListResult{List: v})
}

// Response pagination data object
func (ctx *Context) ResPage(v interface{}, pr *PaginationResult) {
	list := ListResult{
		List:       v,
		Pagination: pr,
	}
	ctx.JSONSuccess(list)
}
