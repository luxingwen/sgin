package model

// 定义权限位
const (
	Create = 1 << 3 // 创建权限: 1000 (8)
	Read   = 1 << 2 // 查询权限: 100 (4)
	Edit   = 1 << 1 // 编辑权限: 010 (2)
	Delete = 1 << 0 // 删除权限: 001 (1)
)

// Permission 定义了权限的基础信息
type Permission struct {
	Id         uint   `gorm:"primary_key" json:"id"`                  // ID 是权限的主键
	Uuid       string `gorm:"type:char(36);primary_key" json:"uuid"`  // UUID 是权限的唯一标识符
	Name       string `gorm:"type:varchar(50);index" json:"name"`     // Name 是权限的名称
	Bit        uint   `gorm:"type:int(3)" json:"bit"`                 // Bit 是权限的位
	ParentUuid string `gorm:"type:char(36);index" json:"parent_uuid"` // ParentUuid 是权限的父级 UUID
	CreatedAt  string `gorm:"autoCreateTime" json:"created_at"`       // CreatedAt 记录了权限创建的时间
	UpdatedAt  string `gorm:"autoUpdateTime" json:"updated_at"`       // UpdatedAt 记录了权限信息最后更新的时间
}

type PermissionMenu struct {
	Id             uint   `gorm:"primary_key" json:"id"`                      // ID 是权限菜单的主键
	Uuid           string `gorm:"type:char(36);primary_key" json:"uuid"`      // UUID 是权限菜单的唯一标识符
	PermissionUuid string `gorm:"type:char(36);index" json:"permission_uuid"` // PermissionUuid 是权限的 UUID
	MenuUuid       string `gorm:"type:char(36);index" json:"menu_uuid"`       // MenuUuid 是菜单的 UUID
	CreatedAt      string `gorm:"autoCreateTime" json:"-"`                    // CreatedAt 记录了权限菜单创建的时间
	UpdatedAt      string `gorm:"autoUpdateTime" json:"-"`                    // UpdatedAt 记录了权限菜单信息最后更新的时间
}

type ReqPermissionMenuCreate struct {
	PermissionUuid string   `json:"permission_uuid" binding:"required"` // PermissionUuid 是权限的 UUID
	MenuUuids      []string `json:"menu_uuids" binding:"required"`      // MenuUuids 是菜单的 UUID 列表

}

// 用户权限关联
type UserPermission struct {
	Id             uint   `gorm:"primary_key" json:"id"`                      // ID 是用户权限关联的主键
	Uuid           string `gorm:"type:char(36);primary_key" json:"uuid"`      // UUID 是用户权限关联的唯一标识符
	UserUuid       string `gorm:"type:char(36);index" json:"user_uuid"`       // UserUuid 是用户的 UUID
	PermissionUuid string `gorm:"type:char(36);index" json:"permission_uuid"` // PermissionUuid 是权限的 UUID
	CreatedAt      string `gorm:"autoCreateTime" json:"created_at"`           // CreatedAt 记录了用户权限关联创建的时间
	UpdatedAt      string `gorm:"autoUpdateTime" json:"updated_at"`           // UpdatedAt 记录了用户权限关联信息最后更新的时间
}

type ReqPermissionUserCreate struct {
	UserUuid        string   `json:"user_uuid" binding:"required"`        // UserUuid 是用户的 UUID
	PermissionUuids []string `json:"permission_uuids" binding:"required"` // PermissionUuids 是权限的 UUID 列表
}
