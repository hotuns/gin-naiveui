package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Enable   bool   `json:"enable"`
}

func (User) TableName() string {
	return "user"
}
