package inout

import (
	"gin-naiveui/model"
	"time"
)

type LoginRes struct {
	AccessToken string `json:"accessToken"`
}

type UserDetailRes struct {
	ID       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Enable   bool   `json:"enable" binding:"required"`

	Profile     *model.Profile `json:"profile"`
	Roles       []*model.Role  `json:"roles"`
	CurrentRole *model.Role    `json:"currentRole"`
}

type RoleListRes []*model.Role

type UserListItem struct {
	ID          int           `json:"id"`
	Username    string        `json:"username"`
	Enable      bool          `json:"enable"`
	CreatedTime time.Time     `json:"createTime"`
	UpdatedTime time.Time     `json:"updateTime"`
	Gender      int           `json:"gender"`
	Avatar      string        `json:"avatar"`
	Address     string        `json:"address"`
	Email       string        `json:"email"`
	Roles       []*model.Role `json:"roles"`
}

type UserListRes struct {
	PageData []UserListItem `json:"pageData"`
	Total    int64          `json:"total"`
}

type RoleListPageItem struct {
	model.Role
	PermissionIds []int64 `json:"permissionIds" gorm:"-"`
}

type RoleListPageRes struct {
	PageData []RoleListPageItem `json:"pageData"`
	Total    int64              `json:"total"`
}
