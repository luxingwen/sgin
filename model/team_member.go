package model

import "time"

// TeamMember 定义了团队成员的基础信息
type TeamMember struct {
	Id        uint      `gorm:"primary_key" json:"id"`            // ID 是团队成员的主键
	UUID      string    `gorm:"type:char(36);unique" json:"uuid"` // UUID 是团队成员的唯一标识符
	TeamUUID  string    `gorm:"type:char(36)" json:"team_uuid"`   // TeamUUID 是团队的UUID
	UserUUID  string    `gorm:"type:char(36)" json:"user_uuid"`   // UserUUID 是用户的UUID
	Role      string    `gorm:"type:varchar(100)" json:"role"`    // Role 是成员在团队中的角色
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了团队成员加入的时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了团队成员信息最后更新的时间
}
