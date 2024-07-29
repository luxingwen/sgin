package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func MigrateDbTable(db *gorm.DB) {
	db.AutoMigrate(
		&AppPermission{},
		&API{},
		&App{},
		&Log{},
		&Menu{},
		&Role{},
		&User{},
		&RoleMenuPermission{},
		&UserRole{},
		&Team{},
		&TeamMember{},
		&VerificationCode{},
		&SysLoginLog{},
		&SysOpLog{},
		&SysAPI{},
	)

	// 创建默认用户
	var user User
	// 查询用户是否存在
	err := db.Where("username = ?", "admin").First(&user).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		// 创建默认用户
		user = User{
			Uuid:      uuid.New().String(),
			Username:  "admin",
			Password:  "jV+K7rZOPOILU30ExIZAfq9IlkZhfPz0k+dvW3lPoIA=",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		err = db.Create(&user).Error
		if err != nil {
			log.Fatal("Failed to create default user", err)
		}
	}

}
