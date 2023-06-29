package model

type PagedResponse struct {
	Data     interface{} `json:"data"`
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
}
