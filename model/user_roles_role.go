package model

type UserRolesRole struct {
	UserId int `gorm:"column:user_id" json:"userId"`
	RoleId int `gorm:"column:role_id" json:"roleId"`
}

func (UserRolesRole) TableName() string {
	return "user_roles_role"
}
