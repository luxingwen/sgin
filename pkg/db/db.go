package db

import (
	"fmt"
	"log"
	"sgin/pkg/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDB(cfg config.MySQLConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	// avoid leaking real password in logs
	safeDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, "*****", cfg.Host, cfg.Port, cfg.Database)
	log.Println("dsn:", safeDSN)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to mysql: %v", err)
	}
	// 配置数据库连接池，提升稳定性
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(1 * time.Hour)
	}
	return db
}
