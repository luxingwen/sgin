package model

import "gorm.io/gorm"

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
	)
}
