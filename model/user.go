package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreateTime time.Time      `json:"createTime"`
	UpdateTime time.Time      `json:"updateTime"`
	DeleteTime gorm.DeletedAt `gorm:"index" json:"deleteTime"`

	Username string `json:"username"`
	Password string `json:"password"`
	Enable   bool   `json:"enable"`
}

func (User) TableName() string {
	return "user"
}
