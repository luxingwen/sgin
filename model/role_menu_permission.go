package model

import "time"

// RoleMenuPermission 定义了角色菜单权限的基础信息
type RoleMenuPermission struct {
	Id        uint      `gorm:"primary_key" json:"id"`                // ID 是角色菜单权限的主键
	UUID      string    `gorm:"type:char(36);index" json:"uuid"`      // UUID 是角色菜单权限的唯一标识符
	RoleUUID  string    `gorm:"type:char(36);index" json:"role_uuid"` // RoleUUID 是角色的 UUID
	MenuUUID  string    `gorm:"type:char(36);index" json:"menu_uuid"` // MenuUUID 是菜单的 UUID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`     // CreatedAt 记录了角色菜单权限创建的时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`     // UpdatedAt 记录了角色菜单权限信息最后更新的时间
}
