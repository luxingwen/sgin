package model

type Pagination struct {
	PageSize  int    `form:"pageSize" json:"pageSize"`
	Current   int    `form:"current" json:"current"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func (p *Pagination) GetOffset() int {
	return p.PageSize * (p.Current - 1)
}

// 角色查询参数
type ReqRoleQueryParam struct {
	Name     string `form:"name"`
	IsActive bool   `form:"is_active"`
	Pagination
}

type ReqUserLogin struct {
	// 用户名或邮箱
	Username string `json:"username" binding:"required"`
	// 密码
	Password string `json:"password" binding:"required"`
}

type ReqUserQueryParam struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Sex      int    `json:"sex"`
	Username string `json:"username"`
	Status   int    `json:"status"`
	Uuid     string `json:"uuid"`
	Pagination
}

// 菜单查询参数
type ReqMenuQueryParam struct {
	Name string `json:"name"`
	Pagination
}

// 删除用户删除参数
type ReqUserDeleteParam struct {
	Uuid string `json:"uuid" binding:"required"`
}

// 删除菜单参数
type ReqMenuDeleteParam struct {
	Uuid string `json:"uuid" binding:"required"`
}

type ReqApiQueryParam struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	Method string `json:"method"`
	Status int    `json:"status"`
	Pagination
}

// 查询app的参数
type ReqAppQueryParam struct {
	Name   string `json:"name"`
	ApiKey string `json:"api_key"`
	Status int    `json:"status"`
	Pagination
}

// uuid参数
type ReqUuidParam struct {
	Uuid  string   `json:"uuid"`
	Uuids []string `json:"uuids"`
}

// api权限参数
type ReqApiPermissionParam struct {
	AppId string `json:"app_id"`
	Pagination
}
