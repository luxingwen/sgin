package model

import "time"

// UserRole 定义了用户角色的基础信息
type UserRole struct {
	UUID      string    `gorm:"type:char(36);primary_key" json:"uuid"` // UUID 是用户角色的唯一标识符
	UserUUID  string    `gorm:"type:char(36)" json:"user_uuid"`        // UserUUID 是用户的UUID
	RoleUUID  string    `gorm:"type:char(36)" json:"role_uuid"`        // RoleUUID 是角色的UUID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`      // CreatedAt 记录了用户获得角色的时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`      // UpdatedAt 记录了用户角色信息最后更新的时间
}
