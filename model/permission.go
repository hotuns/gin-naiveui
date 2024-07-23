package model

type Permission struct {
	ID          int          `json:"id" gorm:"primaryKey" `
	Name        string       `json:"name"`
	Code        string       `json:"code"`
	Type        string       `json:"type"`
	ParentId    *int         `json:"parentId" gorm:"column:parent_id"`
	Path        string       `json:"path"`
	Redirect    string       `json:"redirect"`
	Icon        string       `json:"icon"`
	Component   string       `json:"component"`
	Layout      string       `json:"layout"`
	KeepAlive   bool         `json:"keepAlive" column:"keep_alive"`
	Method      string       `json:"method"`
	Description string       `json:"description"`
	Show        bool         `json:"show"`
	Enable      bool         `json:"enable"`
	SortOrder   int          `json:"sortOrder" gorm:"column:sort_order"`
	Children    []Permission `json:"children" gorm:"-"`
}

func (Permission) TableName() string {
	return "permission"
}
