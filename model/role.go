package model

import "time"

// Role 结构体定义了角色的基础信息
type Role struct {
	ID        uint      `gorm:"primary_key" json:"id"`                // ID 是角色的主键
	Uuid      string    `gorm:"type:varchar(100);unique" json:"uuid"` // Uuid 是角色的唯一标识
	TeamUuid  string    `gorm:"type:char(36)" json:"team_uuid"`       // TeamUuid 是团队的UUID
	Name      string    `gorm:"type:varchar(100);unique" json:"name"` // Name 是角色的名称，它在系统中是唯一的
	Desc      string    `gorm:"type:varchar(255)" json:"desc"`        // Desc 是对角色的描述
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`     // CreatedAt 记录了角色创建的时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`     // UpdatedAt 记录了角色最后更新的时间
	IsActive  bool      `gorm:"default:true" json:"is_active"`        // IsActive 标识角色是否是活跃的
}
