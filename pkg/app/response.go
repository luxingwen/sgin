package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	return
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
