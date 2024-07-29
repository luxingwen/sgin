package model

import "time"

// Menu 定义了菜单的基础信息
type Menu struct {
	Id         uint      `gorm:"primary_key" json:"id"`                      // ID 是菜单的主键
	UUID       string    `gorm:"type:char(36);index" json:"uuid"`            // UUID 是菜单的唯一标识符
	Name       string    `gorm:"type:varchar(100)" json:"name"`              // Name 是菜单的名称
	Link       string    `gorm:"type:varchar(255)" json:"link"`              // Link 是菜单的链接
	ParentUUID string    `gorm:"type:char(36)" json:"parent_uuid"`           // ParentUUID 是父菜单的UUID
	Icon       string    `gorm:"type:varchar(100)" json:"icon"`              // Icon 是菜单的图标
	Order      int       `json:"order" gorm:"column:order;type:int"`         // Order 是菜单的排序
	IsShow     bool      `json:"is_show" gorm:"column:is_show;type:boolean"` // IsShow 是菜单是否显示
	Type       int       `json:"type" gorm:"column:type;type:int"`           // Type 是菜单的类型  1:目录 2、菜单、 3:按钮  4:链接
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"-"`                    // CreatedAt 记录了菜单创建的时间
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"-"`                    // UpdatedAt 记录了菜单信息最后更新的时间
}

// MenuAPI 结构体用于表示关联关系
type MenuAPI struct {
	Id        uint      `gorm:"primary_key" json:"id"` // ID 是关联关系的主键
	MenuUUID  string    `gorm:"type:char(36);index" json:"menu_uuid"`
	APIUUID   string    `gorm:"type:char(36);index" json:"api_uuid"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
