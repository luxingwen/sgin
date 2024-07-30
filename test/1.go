package main

import (
	"fmt"
)

// 定义权限位
const (
	Read   = 1 << 2 // 查询权限: 100 (4)
	Create = 1 << 3 // 创建权限: 1000 (8)
	Edit   = 1 << 1 // 编辑权限: 010 (2)
	Delete = 1 << 0 // 删除权限: 001 (1)
)

// 用户结构体定义
type User struct {
	ID       int    // 用户ID
	Name     string // 用户名称
	Email    string // 用户邮箱
	RoleBits int    // 用户的权限位
}

// 设置权限
func (u *User) SetPermission(permission int) {
	u.RoleBits |= permission
}

// 检查用户是否有指定权限
func (u *User) HasPermission(permission int) bool {
	return u.RoleBits&permission == permission
}

func main() {
	// 创建用户
	user := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	// 赋予用户查询和创建权限
	user.SetPermission(Read | Create)

	fmt.Println("User:", user.RoleBits)

	// 检查用户是否有各类权限
	fmt.Printf("User has read permission: %v\n", user.HasPermission(Read))
	fmt.Printf("User has create permission: %v\n", user.HasPermission(Create))
	fmt.Printf("User has edit permission: %v\n", user.HasPermission(Edit))
	fmt.Printf("User has delete permission: %v\n", user.HasPermission(Delete))

	// 赋予用户查询和创建权限
	user.SetPermission(Read | Create | Edit | Delete)

	fmt.Println("User:", user.RoleBits)

	// 检查用户是否有各类权限
	fmt.Printf("User has read permission: %v\n", user.HasPermission(Read))
	fmt.Printf("User has create permission: %v\n", user.HasPermission(Create))
	fmt.Printf("User has edit permission: %v\n", user.HasPermission(Edit))
	fmt.Printf("User has delete permission: %v\n", user.HasPermission(Delete))
}
