package db

import (
	"fmt"
	"log"
	"time"

	"github.com/luxingwen/sgin/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB(dbType string, cfg config.DBConfig) *gorm.DB {
	var dialector gorm.Dialector
	var dsn string

	if dbType == "postgres" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port)
		dialector = postgres.Open(dsn)
	} else {
		// Default to MySQL
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		dialector = mysql.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", dbType, err)
	}

	// 配置数据库连接池
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(1 * time.Hour)
	}
	return db
}
