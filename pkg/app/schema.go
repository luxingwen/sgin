package app

type ListResult struct {
	List       interface{}       `json:"list"`
	Pagination *PaginationResult `json:"pagination,omitempty"`
}

type PaginationResult struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}

// PaginationRequest 通用分页请求结构，可在各 handler 中复用或组合使用。
// 前端发送 {"page":1,"per_page":20} 格式的 JSON，或者使用 query params。
type PaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// Normalize 返回规范化的分页参数：current page、pageSize、offset
func (p *PaginationRequest) Normalize() (current int, pageSize int, offset int) {
	current = p.Page
	pageSize = p.PerPage
	if current <= 0 {
		current = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset = (current - 1) * pageSize
	return
}

// Filter 定义单个过滤条件。前端通过数组传入多个条件。
// 支持的 op（示例）："eq","ne","like","in","gt","gte","lt","lte","between"
type Filter struct {
	Field string      `json:"field"` // 字段名（服务端应校验白名单）
	Op    string      `json:"op"`    // 操作符
	Value interface{} `json:"value"` // 值或值数组（in/between）
}

// Sorts 为排序字段数组，支持前缀 - 表示降序，例如 ["-created_at","name"]
type Sorts []string
