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
	Pagination
}
