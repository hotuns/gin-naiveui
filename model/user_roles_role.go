package model

type UserRolesRole struct {
	UserId int `gorm:"column:user_id"`
	RoleId int `gorm:"column:role_id"`
}

func (UserRolesRole) TableName() string {
	return "user_roles_role"
}
