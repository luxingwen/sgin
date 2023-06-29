package model

import "time"

// Menu 定义了菜单的基础信息
type Menu struct {
	Id         uint      `gorm:"primary_key" json:"id"`            // ID 是菜单的主键
	UUID       string    `gorm:"type:char(36);index" json:"uuid"`  // UUID 是菜单的唯一标识符
	Name       string    `gorm:"type:varchar(100)" json:"name"`    // Name 是菜单的名称
	Link       string    `gorm:"type:varchar(255)" json:"link"`    // Link 是菜单的链接
	ParentUUID string    `gorm:"type:char(36)" json:"parent_uuid"` // ParentUUID 是父菜单的UUID
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了菜单创建的时间
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了菜单信息最后更新的时间
}
