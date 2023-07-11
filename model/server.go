package model

import "time"

type Server struct {
	Id          int64     `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	UUID        string    `json:"uuid" gorm:"column:uuid;type:varchar(36);not null;unique_index"`
	Name        string    `json:"name" gorm:"column:name;type:varchar(255);not null;unique_index"`
	Address     string    `json:"address" gorm:"column:address;type:varchar(255);not null;"`
	Description string    `json:"description" gorm:"column:description;type:varchar(255);not null;"`
	CreateAt    time.Time `json:"create_at" gorm:"column:create_at;type:datetime;not null;"`
	UpdateAt    time.Time `json:"update_at" gorm:"column:update_at;type:datetime;not null;"`
}
