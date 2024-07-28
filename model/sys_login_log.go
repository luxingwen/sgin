package model

const (
	// 登录状态
	LoginStatusSuccess = 1
	LoginStatusFail    = 2
)

// 系统登录日志
type SysLoginLog struct {
	ID uint `json:"id" gorm:"primaryKey;comment:'主键ID'"`

	Ip string `json:"ip" gorm:"comment:'IP地址'"`

	Username string `json:"username" gorm:"comment:'用户名'"`

	UserAgent string `json:"user_agent" gorm:"comment:'UserAgent'"`
	// 登录状态
	Status int `json:"status" gorm:"comment:'登录状态'"` // 登录状态 1:成功 2:失败
	// 消息
	Message string `json:"message" gorm:"comment:'消息'"` // 消息
	// 浏览器
	Browser string `json:"browser" gorm:"comment:'浏览器'"` // 浏览器
	// 操作系统
	Os string `json:"os" gorm:"comment:'操作系统'"` // 操作系统

	// 登录设备
	Device string `json:"device" gorm:"comment:'登录设备'"` //

	// 地址
	Address string `json:"address" gorm:"comment:'地址'"` // 地址
	// 创建时间
	CreatedAt string `json:"created_at" gorm:"autoCreateTime;comment:'创建时间'"` // 创建时间
}
