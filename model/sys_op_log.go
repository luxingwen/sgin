package model

type SysOpLog struct {
	ID        uint   `json:"id" gorm:"primaryKey;comment:'主键ID'"`                     // 主键ID
	RequestId string `json:"request_id" gorm:"type:varchar(50);index;comment:'请求ID'"` // 请求ID
	UserUuid  string `json:"user_uuid" gorm:"type:char(36);index;comment:'用户UUID'"`   // 用户UUID
	Path      string `json:"path" gorm:"comment:'请求路径'"`                              // 请求路径
	Method    string `json:"method" gorm:"comment:'请求方法'"`                            // 请求方法
	Ip        string `json:"ip" gorm:"comment:'请求IP'"`                                // 请求IP
	Status    int    `json:"status" gorm:"comment:'状态'"`                              // 状态
	Code      int    `json:"code" gorm:"comment:'状态码'"`                               // 状态码
	Message   string `json:"message" gorm:"comment:'消息'"`                             // 消息
	Params    string `json:"params" gorm:"type:text;comment:'请求参数'"`                  // 请求参数
	Duration  int64  `json:"duration" gorm:"comment:'请求耗时(毫秒)'"`                      // 请求耗时
	CreatedAt string `json:"created_at" gorm:"autoCreateTime;comment:'创建时间'"`         // 创建时间
}

type SysOpLogRes struct {
	SysOpLog
	Username string `json:"username" gorm:"comment:'用户名'"` // 用户名
	Module   string `json:"module" gorm:"comment:'模块'"`    // 模块
	Name     string `json:"name" gorm:"comment:'操作名称'"`    // 操作名称
}
