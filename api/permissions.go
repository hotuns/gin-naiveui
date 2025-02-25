package api

import (
	"gin-naiveui/db"
	"gin-naiveui/inout"
	"gin-naiveui/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Permissions = &permissions{}

type permissions struct {
}

func (permissions) List(c *gin.Context) {
	var onePermissList = make([]model.Permission, 0)
	// 一级菜单
	db.Dao.Model(model.Permission{}).Where("parent_id is NULL").Order("sort_order Asc").Find(&onePermissList)

	for i, perm := range onePermissList {
		onePermissList[i].Children = getChildren(uint(perm.ID))
	}
	Resp.Succ(c, onePermissList)
}

func getChildren(parentID uint) []model.Permission {
	var children []model.Permission
	db.Dao.Model(model.Permission{}).Where("parent_id = ?", parentID).Order("sort_order Asc").Find(&children)

	for i, child := range children {
		children[i].Children = getChildren(uint(child.ID))
	}

	return children
}

func (permissions) ListPage(c *gin.Context) {
	var data = &inout.RoleListPageRes{}
	var name = c.DefaultQuery("name", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	orm := db.Dao.Model(model.Role{})
	if name != "" {
		orm = orm.Where("name like ?", "%"+name+"%")
	}
	orm.Count(&data.Total)

	orm.Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&data.PageData)
	for i, datum := range data.PageData {
		var perIdList []int64
		db.Dao.Model(model.RolePermissionsPermission{}).Where("role_id=?", datum.ID).Select("permission_id").Find(&perIdList)
		data.PageData[i].PermissionIds = perIdList
	}
	Resp.Succ(c, data)
}
func (permissions) Add(c *gin.Context) {
	var params inout.AddPermissionReq
	err := c.Bind(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}

	err = db.Dao.Model(model.Permission{}).Create(&model.Permission{
		Name:      params.Name,
		Code:      params.Code,
		Type:      params.Type,
		ParentId:  params.ParentId, // insert value null
		Path:      params.Path,
		Icon:      params.Icon,
		Component: params.Component,
		Layout:    params.Layout,
		KeepAlive: (params.KeepAlive),
		Show:      (params.Show),
		Enable:    (params.Enable),
		SortOrder: params.SortOrder,
	}).Error
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}

	Resp.Succ(c, "")
}
func (permissions) Delete(c *gin.Context) {
	id := c.Param("id")
	err := db.Dao.Transaction(func(tx *gorm.DB) error {
		tx.Where("id =?", id).Delete(&model.Permission{})
		tx.Where("permission_id =?", id).Delete(&model.RolePermissionsPermission{})
		return nil
	})
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}
func (permissions) PatchPermission(c *gin.Context) {
	var params inout.PatchPermissionReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	err = db.Dao.Model(model.Permission{}).Where("id=?", params.Id).Updates(model.Permission{
		Name:      params.Name,
		Code:      params.Code,
		Type:      params.Type,
		ParentId:  params.ParentId,
		Path:      params.Path,
		Icon:      params.Icon,
		Component: params.Component,
		Layout:    params.Layout,
		KeepAlive: params.KeepAlive,
		Method:    params.Component,
		Show:      params.Show,
		Enable:    params.Enable,
		SortOrder: params.SortOrder,
	}).Error
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}

	if !params.Show {
		db.Dao.Model(model.Permission{}).Where("id=?", params.Id).Update("show", false)
	}

	Resp.Succ(c, "")

}
func (permissions) ValidateMenuPath(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		Resp.Err(c, 20001, "path is required")
		return
	}

	var menus []model.Permission
	// 获取所有类型为MENU的权限记录
	err := db.Dao.Model(model.Permission{}).Where("type = ?", "MENU").Find(&menus).Error
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}

	// 检查是否存在匹配的菜单路径
	hasMenu := false
	for _, menu := range menus {
		if menu.Path == path {
			hasMenu = true
			break
		}
	}

	Resp.Succ(c, hasMenu)
}
