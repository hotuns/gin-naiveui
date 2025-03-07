package api

import (
	"crypto/md5"
	"fmt"
	"gin-naiveui/db"
	"gin-naiveui/inout"
	"gin-naiveui/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var User = &user{}

type user struct {
}

func (user) Detail(c *gin.Context) {
	var data = &inout.UserDetailRes{}

	var uid, _ = c.Get("uid")
	db.Dao.Model(model.User{}).Where("id=?", uid).Find(&data)
	fmt.Println(data)

	db.Dao.Model(model.Profile{}).Where("user_id=?", uid).Find(&data.Profile)

	urolIdList := db.Dao.Model(model.UserRolesRole{}).Where("user_id=?", uid).Select("role_id")
	db.Dao.Model(model.Role{}).Where("id IN (?)", urolIdList).Find(&data.Roles)
	if len(data.Roles) > 0 {
		data.CurrentRole = data.Roles[0]
	}
	Resp.Succ(c, data)
}

func (user) List(c *gin.Context) {
	var data = inout.UserListRes{
		PageData: make([]inout.UserListItem, 0),
	}
	var gender = c.DefaultQuery("gender", "")
	var enable = c.DefaultQuery("enable", "")
	var username = c.DefaultQuery("username", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	var profileList []model.Profile
	orm := db.Dao.Model(model.Profile{})
	if gender != "" {
		orm = orm.Where("gender=?", gender)
	}
	if enable != "" {
		orm = orm.Where("user_id in(?)", db.Dao.Model(model.User{}).Where("enable=?", enable).Select("id"))
	}
	if username != "" {
		orm = orm.Where("nick_name like ?", "%"+username+"%")
	}

	orm.Count(&data.Total)
	orm.Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&profileList)
	for _, datum := range profileList {
		var uinfo model.User
		db.Dao.Model(model.User{}).Where("id=?", datum.UserId).First(&uinfo)
		var rols []*model.Role
		db.Dao.Model(model.Role{}).Where("id IN (?)", db.Dao.Model(model.UserRolesRole{}).Where("user_id=?", datum.UserId).Select("role_id")).Find(&rols)
		data.PageData = append(data.PageData, inout.UserListItem{
			ID:          int(uinfo.ID),
			Username:    uinfo.Username,
			Enable:      uinfo.Enable,
			CreatedTime: uinfo.CreateTime,
			UpdatedTime: uinfo.UpdateTime,
			Gender:      datum.Gender,
			Avatar:      datum.Avatar,
			Address:     datum.Address,
			Email:       datum.Email,
			Roles:       rols,
		})
	}
	Resp.Succ(c, data)
}

func (user) Profile(c *gin.Context) {
	var params inout.PatchProfileUserReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	err = db.Dao.Model(model.Profile{}).Where("id=?", params.Id).Updates(model.Profile{
		Avatar:   params.Avatar,
		Gender:   params.Gender,
		Address:  params.Address,
		Email:    params.Email,
		NickName: params.NickName,
	}).Error

	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}
func (user) Update(c *gin.Context) {
	var params inout.PatchUserReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	updates := map[string]interface{}{}
	if params.Password != nil {
		updates["password"] = fmt.Sprintf("%x", md5.Sum([]byte(*params.Password)))
	}
	if params.Enable != nil {
		updates["enable"] = *params.Enable
	}
	if params.Username != nil {
		updates["username"] = *params.Username
		db.Dao.Model(model.Profile{}).Where("user_id=?", params.Id).Update("nickName", *params.Username)
	}
	if len(updates) > 0 {
		db.Dao.Model(model.User{}).Where("id=?", params.Id).Updates(updates)
	}
	if params.RoleIds != nil {
		db.Dao.Where("user_id=?", params.Id).Delete(model.UserRolesRole{})
		if len(*params.RoleIds) > 0 {
			for _, i2 := range *params.RoleIds {
				db.Dao.Model(model.UserRolesRole{}).Create(&model.UserRolesRole{
					UserId: params.Id,
					RoleId: i2,
				})
			}
		}
	}
	Resp.Succ(c, err)
}

func (user) Add(c *gin.Context) {
	var params inout.AddUserReq
	err := c.Bind(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	err = db.Dao.Transaction(func(tx *gorm.DB) error {
		var newUser = model.User{
			Username: params.Username,
			Password: fmt.Sprintf("%x", md5.Sum([]byte(params.Password))),
			Enable:   params.Enable,
		}
		err = tx.Create(&newUser).Error
		if err != nil {
			return err
		}
		tx.Create(&model.Profile{
			UserId:   int(newUser.ID),
			NickName: newUser.Username,
		})
		for _, id := range params.RoleIds {
			tx.Create(&model.UserRolesRole{
				UserId: int(newUser.ID),
				RoleId: id,
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
func (user) Delete(c *gin.Context) {
	uid := c.Param("id")
	err := db.Dao.Transaction(func(tx *gorm.DB) error {
		tx.Where("id =?", uid).Delete(&model.User{})
		tx.Where("user_id =?", uid).Delete(&model.UserRolesRole{})
		tx.Where("user_id =?", uid).Delete(&model.Profile{})
		return nil
	})
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}
