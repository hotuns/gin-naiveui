package model

type RolePermissionsPermission struct {
	RoleId       int `gorm:"column:role_id"`
	PermissionId int `gorm:"column:permission_id"`
}

func (RolePermissionsPermission) TableName() string {
	return "role_permissions_permission"
}
