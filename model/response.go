package model

type PagedResponse struct {
	Data     interface{} `json:"data"`
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
}

type ResUserLogin struct {
	Token string `json:"token"`
}

type BaseResponse struct {
	TraceID string `json:"trace_id"` // 请求唯一标识
	Code    int    `json:"code"`     // 状态码
	Message string `json:"message"`  // 提示信息
}

type UserInfoResponse struct {
	BaseResponse
	Data User `json:"data"`
}

type BasePageResponse struct {
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

type UserPageResponse struct {
	BasePageResponse
	Data []User `json:"data"`
}

type UserQueryResponse struct {
	BaseResponse
	Data UserPageResponse `json:"data"`
}

type AppInfoResponse struct {
	BaseResponse
	Data App `json:"data"`
}

type AppPageResponse struct {
	BasePageResponse
	Data []App `json:"data"`
}

type AppQueryResponse struct {
	BaseResponse
	Data AppPageResponse `json:"data"`
}

type MenuInfoResponse struct {
	BaseResponse
	Data Menu `json:"data"`
}

type MenuPageResponse struct {
	BasePageResponse
	Data []Menu `json:"data"`
}

type MenuQueryResponse struct {
	BaseResponse
	Data MenuPageResponse `json:"data"`
}

type RoleInfoResponse struct {
	BaseResponse
	Data Role `json:"data"`
}

type RolePageResponse struct {
	BasePageResponse
	Data []Role `json:"data"`
}

type RoleQueryResponse struct {
	BaseResponse
	Data RolePageResponse `json:"data"`
}

type ServerInfoResponse struct {
	BaseResponse
	Data Server `json:"data"`
}

type ServerPageResponse struct {
	BasePageResponse
	Data []Server `json:"data"`
}

type ServerQueryResponse struct {
	BaseResponse
	Data ServerPageResponse `json:"data"`
}

type TeamInfoResponse struct {
	BaseResponse
	Data Team `json:"data"`
}

type TeamPageResponse struct {
	BasePageResponse
	Data []Team `json:"data"`
}

type TeamQueryResponse struct {
	BaseResponse
	Data TeamPageResponse `json:"data"`
}

type TeamMemberInfoResponse struct {
	BaseResponse
	Data TeamMember `json:"data"`
}
