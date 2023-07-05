package model

import "time"

// 验证码
type VerificationCode struct {
	Id        uint      `gorm:"primary_key" json:"id"`            // ID 是验证码的主键
	UUID      string    `gorm:"type:char(36);index" json:"uuid"`  // UUID 是验证码的唯一标识符
	Code      string    `gorm:"type:varchar(6)" json:"code"`      // Code 是验证码的内容
	Email     string    `gorm:"type:varchar(100)" json:"email"`   // Email 是验证码的接收者
	Phone     string    `gorm:"type:varchar(11)" json:"phone"`    // Phone 是验证码的接收者
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了验证码创建的时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了验证码信息最后更新的时间
	Status    int       `gorm:"type:int(1)" json:"status"`        // Status 0:未使用 1:已使用 2:已过期
}
