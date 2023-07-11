package model

import "time"

// Log 定义了用户操作日志的基础信息
type Log struct {
	Id        uint      `gorm:"primary_key" json:"id"`                // ID 是日志的主键
	UUID      string    `gorm:"type:char(36);index" json:"uuid"`      // UUID 是日志的唯一标识符
	UserUUID  string    `gorm:"type:char(36);index" json:"user_uuid"` // UserUUID 是用户的 UUID
	Action    string    `gorm:"type:varchar(255)" json:"action"`      // Action 是用户的操作内容
	AppId     string    `gorm:"type:char(36);index" json:"app_id"`    // AppId 是应用的 UUID
	ReqBody   string    `gorm:"type:text" json:"req_body"`            // ReqBody 是用户的请求内容
	RespBody  string    `gorm:"type:text" json:"resp_body"`           // RespBody 是用户的响应内容
	Status    int       `gorm:"type:int" json:"status"`               // Status 是用户的操作状态
	Method    string    `gorm:"type:varchar(10)" json:"method"`       // Method 是用户的HTTP方法
	Path      string    `gorm:"type:varchar(255)" json:"path"`        // Path 是用户的请求路径
	Ip        string    `gorm:"type:varchar(255)" json:"ip"`          // Ip 是用户的IP地址
	Message   string    `gorm:"type:text" json:"message"`             // Message 是用户的操作信息
	Header    string    `gorm:"type:text" json:"header"`              // Header 是用户的请求头
	TraceID   string    `gorm:"type:varchar(255)" json:"trace_id"`    // TraceID 是用户的请求追踪ID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`     // CreatedAt 记录了日志创建的时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`     // UpdatedAt 记录了日志信息最后更新的时间
}
