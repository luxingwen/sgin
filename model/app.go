package model

import "time"

// APP 定义了调用方的基础信息
type App struct {
	Id        uint      `gorm:"primary_key" json:"id"`            // ID 是调用方的主键
	UUID      string    `gorm:"type:char(36);index" json:"uuid"`  // UUID 是调用方的唯一标识符
	Name      string    `gorm:"type:varchar(100)" json:"name"`    // Name 是调用方的名称
	ApiKey    string    `gorm:"type:varchar(255)" json:"api_key"` // ApiKey 是调用方的API Key
	SecKey    string    `gorm:"type:varchar(255)" json:"sec_key"` // SecKey 是调用方的Sec Key
	UserUUID  string    `gorm:"type:char(36)" json:"user_uuid"`   // UserUUID 是用户的UUID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了调用方创建的时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了调用方信息最后更新的时间
	Status    int       `gorm:"type:int(1)" json:"status"`        // Status 0:未启用 1:启用 2:删除
}
