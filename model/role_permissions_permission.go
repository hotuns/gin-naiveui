package model

type RolePermissionsPermission struct {
	RoleId       int `gorm:"column:role_id" json:"roleId"`
	PermissionId int `gorm:"column:permission_id" json:"permissionId"`
}

func (RolePermissionsPermission) TableName() string {
	return "role_permissions_permission"
}
