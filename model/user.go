package model

import "time"

type User struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Uuid      string    `gorm:"type:char(36);unique" json:"uuid"`         // 用户唯一标识
	Email     string    `gorm:"type:varchar(100);unique" json:"email"`    // 邮箱
	Username  string    `gorm:"type:varchar(100);unique" json:"username"` // 用户名
	Password  string    `gorm:"type:varchar(100)" json:"password"`        // 密码
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`            // 手机号
	Avatar    string    `gorm:"type:varchar(200)" json:"avatar"`          // 头像
	Nickname  string    `gorm:"type:varchar(50)" json:"nickname"`         // 昵称
	Status    int       `gorm:"type:int" json:"status"`                   // 状态 0:禁用 1:启用 2:删除
	Age       int       `gorm:"type:int" json:"age"`                      // 年龄
	Sex       string    `gorm:"type:varchar(10)" json:"sex"`              // 性别 0:未知 1:男 2:女
	Signed    string    `gorm:"type:varchar(255)" json:"signed"`          // 个性签名
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`         // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`         // 更新时间
}
