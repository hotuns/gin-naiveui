package api

import (
	"gin-naiveui/db"
	"gin-naiveui/inout"
	"gin-naiveui/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Role = &role{}

type role struct {
}

func (role) PermissionsTree(c *gin.Context) {
	var uid, _ = c.Get("uid")

	var adminRole int64
	db.Dao.Model(model.UserRolesRole{}).Where("user_id=? and role_id=1", uid).Count(&adminRole)
	orm := db.Dao.Model(model.Permission{}).Where("parent_id is NULL").Order("sort_order Asc")

	if adminRole == 0 {
		uroleIdList := db.Dao.Model(model.UserRolesRole{}).Where("user_id=?", uid).Select("role_id")
		rpermisId := db.Dao.Model(model.RolePermissionsPermission{}).Where("role_id in(?)", uroleIdList).Select("permission_id")
		orm = orm.Where("id in(?)", rpermisId)
	}

	var onePermissList []model.Permission
	orm.Find(&onePermissList)

	for i, perm := range onePermissList {
		onePermissList[i].Children = getChildren(uint(perm.ID))
	}
	Resp.Succ(c, onePermissList)
}

func (role) List(c *gin.Context) {
	var data = &inout.RoleListRes{}

	db.Dao.Model(model.Role{}).Find(&data)

	Resp.Succ(c, data)
}
func (role) ListPage(c *gin.Context) {
	var data = &inout.RoleListPageRes{}
	var name = c.DefaultQuery("name", "")
	var enable = c.DefaultQuery("enable", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	orm := db.Dao.Model(model.Role{})
	if name != "" {
		orm = orm.Where("name like ?", "%"+name+"%")
	}
	if enable != "" {
		ena := false
		if enable == "1" {
			ena = true
		}
		orm = orm.Where("enable = ?", ena)
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
func (role) Update(c *gin.Context) {
	var params inout.PatchRoleReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	orm := db.Dao.Model(model.Role{}).Where("id=?", params.Id)
	if params.Name != nil {
		orm.Update("name", *params.Name)
	}
	if params.Enable != nil {
		orm.Update("enable", *params.Enable)
	}
	if params.Code != nil {
		orm.Update("code", *params.Code)
	}
	if params.PermissionIds != nil {
		db.Dao.Where("role_id=?", params.Id).Delete(model.RolePermissionsPermission{})
		if len(*params.PermissionIds) > 0 {
			for _, i2 := range *params.PermissionIds {
				db.Dao.Model(model.RolePermissionsPermission{}).Create(&model.RolePermissionsPermission{
					PermissionId: i2,
					RoleId:       params.Id,
				})
			}
		}
	}
	Resp.Succ(c, err)
}

func (role) Add(c *gin.Context) {
	var params inout.AddRoleReq
	err := c.Bind(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	err = db.Dao.Transaction(func(tx *gorm.DB) error {
		var record = model.Role{
			Code:   params.Code,
			Name:   params.Name,
			Enable: params.Enable,
		}
		err = tx.Create(&record).Error
		if err != nil {
			return err
		}

		for _, id := range params.PermissionIds {
			tx.Create(&model.RolePermissionsPermission{
				RoleId:       record.ID,
				PermissionId: id,
			})
		}
		return nil
	})
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}

func (role) Delete(c *gin.Context) {
	uid := c.Param("id")
	err := db.Dao.Transaction(func(tx *gorm.DB) error {
		tx.Where("id =?", uid).Delete(&model.Role{})
		tx.Where("role_id =?", uid).Delete(&model.UserRolesRole{})
		tx.Where("role_id =?", uid).Delete(&model.RolePermissionsPermission{})
		return nil
	})
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}
func (role) AddUser(c *gin.Context) {
	var params inout.PatchRoleOpeateUserReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	uid, _ := strconv.Atoi(c.Param("id"))
	params.Id = uid
	db.Dao.Where("user_id in (?) and role_id = ?", params.UserIds, params.Id).Delete(model.UserRolesRole{})
	for _, id := range params.UserIds {
		db.Dao.Model(model.UserRolesRole{}).Create(model.UserRolesRole{
			UserId: id,
			RoleId: params.Id,
		})
	}
	Resp.Succ(c, "")
}
func (role) RemoveUser(c *gin.Context) {
	var params inout.PatchRoleOpeateUserReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	uid, _ := strconv.Atoi(c.Param("id"))
	params.Id = uid
	db.Dao.Where("user_id in (?) and role_id = ?", params.UserIds, params.Id).Delete(model.UserRolesRole{})
	Resp.Succ(c, "")
}
