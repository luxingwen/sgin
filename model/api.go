package model

import "time"

// API 定义了API接口的基础信息
type API struct {
	Id         uint      `gorm:"primary_key" json:"id"`                 // ID 是API的主键
	UUID       string    `gorm:"type:char(36);primary_key" json:"uuid"` // UUID 是API的唯一标识符
	ServerUUID string    `gorm:"type:char(36)" json:"server_uuid"`      // ServerUUID 是API所属的服务器的UUID
	Name       string    `gorm:"type:varchar(100)" json:"name"`         // Name 是API的名称
	Path       string    `gorm:"type:varchar(255)" json:"path"`         // Path 是API的路径
	Method     string    `gorm:"type:varchar(10)" json:"method"`        // Method 是API的HTTP方法
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`      // CreatedAt 记录了API创建的时间
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`      // UpdatedAt 记录了API信息最后更新的时间
	Status     int       `gorm:"type:int(1)" json:"status"`             // Status 0:未启用 1:启用 2:删除
}
