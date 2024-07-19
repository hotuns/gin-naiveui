package model

type Role struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

func (Role) TableName() string {
	return "role"
}
