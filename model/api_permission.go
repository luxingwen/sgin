package model

import "time"

// AppPermission 定义了应用接口权限的基础信息
type AppPermission struct {
	Id        uint      `gorm:"primary_key" json:"id"`                 // ID 是应用接口权限的主键
	UUID      string    `gorm:"type:char(36);primary_key" json:"uuid"` // UUID 是应用接口权限的唯一标识符
	AppUUID   string    `gorm:"type:char(36);index" json:"app_uuid"`   // AppUUID 是应用的 UUID
	APIUUID   string    `gorm:"type:char(36);index" json:"api_uuid"`   // APIUUID 是 API 的 UUID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`      // CreatedAt 记录了应用接口权限创建的时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`      // UpdatedAt 记录了应用接口权限信息最后更新的时间
}
