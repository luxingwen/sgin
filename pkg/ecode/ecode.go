package ecode

import "net/http"

// APIError 代表一类对外可见的 API 错误
// Code 当前与 HTTP 状态码一致，后续可演进为独立业务码
type APIError struct {
	HTTPStatus int
	Code       int
	Message    string
}

func (e *APIError) Error() string { return e.Message }

// New 创建一个 APIError，Code 默认为 HTTPStatus
func New(httpStatus int, message string) *APIError {
	return &APIError{HTTPStatus: httpStatus, Code: httpStatus, Message: message}
}

// 便捷构造函数
func BadRequest(msg string) *APIError {
	return New(http.StatusBadRequest, orDefault(msg, "Bad Request"))
}
func Unauthorized(msg string) *APIError {
	return New(http.StatusUnauthorized, orDefault(msg, "Unauthorized"))
}
func Forbidden(msg string) *APIError { return New(http.StatusForbidden, orDefault(msg, "Forbidden")) }
func NotFound(msg string) *APIError  { return New(http.StatusNotFound, orDefault(msg, "Not Found")) }
func Conflict(msg string) *APIError  { return New(http.StatusConflict, orDefault(msg, "Conflict")) }
func InternalError(msg string) *APIError {
	return New(http.StatusInternalServerError, orDefault(msg, "Internal Server Error"))
}
func TooManyRequests(msg string) *APIError {
	return New(http.StatusTooManyRequests, orDefault(msg, "Too Many Requests"))
}
func ServiceUnavailable(msg string) *APIError {
	return New(http.StatusServiceUnavailable, orDefault(msg, "Service Unavailable"))
}

func orDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}
