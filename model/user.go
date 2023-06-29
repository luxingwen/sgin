package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Uuid      string    `json:"uuid"`     // 用户唯一标识
	Email     string    `json:"email"`    // 邮箱
	Username  string    `json:"username"` // 用户名
	Password  string    `json:"password"` // 密码
	Phone     string    `json:"phone"`    // 手机号
	Avatar    string    `json:"avatar"`   // 头像
	Nickname  string    `json:"nickname"` // 昵称
	Status    int       `json:"status"`   // 状态 0:禁用 1:启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
