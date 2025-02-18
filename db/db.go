package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gin-naiveui/config"
	"gin-naiveui/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Dao *gorm.DB

func Init() {
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		},
	)

	port := config.Config("DB_PORT")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_HOST"),
		port,
		config.Config("DB_NAME"),
	)
	fmt.Println(dsn)

	openDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   dbLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("db connection error is %s", err.Error())
	}

	openDb.AutoMigrate(
		&model.User{}, &model.Profile{}, &model.Role{}, &model.Permission{}, &model.UserRolesRole{}, &model.RolePermissionsPermission{},
	)

	dbCon, err := openDb.DB()
	if err != nil {
		log.Fatalf("openDb.DB error is  %s", err.Error())
	}
	dbCon.SetMaxIdleConns(3)
	dbCon.SetMaxOpenConns(10)
	dbCon.SetConnMaxLifetime(time.Hour)
	Dao = openDb

	// 初始化数据
	var count int64
	// 检查 permission 表中是否已有数据
	Dao.Model(&model.Permission{}).Count(&count)

	if count > 0 {
		log.Println("permission 表已有数据，无需初始化")
	} else {
		CreateInitData(Dao)
	}
}

func intPointer(i int) *int {
	return &i
}

func CreateInitData(db *gorm.DB) {

	permissions := []model.Permission{
		{
			Name:      "资源管理",
			Code:      "Resource_Mgt",
			Type:      "MENU",
			ParentId:  intPointer(2),
			Path:      "/pms/resource",
			Icon:      "i-fe:list",
			Component: "/src/views/pms/resource/index.vue",
			Show:      true,
			Enable:    true,
			SortOrder: 1,
		},
		{
			Name:      "系统管理",
			Code:      "SysMgt",
			Type:      "MENU",
			Icon:      "i-fe:grid",
			Show:      true,
			Enable:    true,
			SortOrder: 2,
		},
		{
			Name:      "角色管理",
			Code:      "RoleMgt",
			Type:      "MENU",
			ParentId:  intPointer(2),
			Path:      "/pms/role",
			Icon:      "i-fe:user-check",
			Component: "/src/views/pms/role/index.vue",
			Show:      true,
			Enable:    true,
			SortOrder: 2,
		},
		{
			Name:      "用户管理",
			Code:      "UserMgt",
			Type:      "MENU",
			ParentId:  intPointer(2),
			Path:      "/pms/user",
			Icon:      "i-fe:user",
			Component: "/src/views/pms/user/index.vue",
			KeepAlive: true,
			Show:      true,
			Enable:    true,
			SortOrder: 3,
		},
		{
			Name:      "分配用户",
			Code:      "RoleUser",
			Type:      "MENU",
			ParentId:  intPointer(3),
			Path:      "/pms/role/user/:roleId",
			Icon:      "i-fe:user-plus",
			Component: "/src/views/pms/role/role-user.vue",
			Show:      false,
			Enable:    true,
			SortOrder: 1,
		},
		{
			Name:      "业务示例",
			Code:      "Demo",
			Type:      "MENU",
			Icon:      "i-fe:grid",
			Show:      true,
			Enable:    true,
			SortOrder: 1,
		},
		{
			Name:      "图片上传",
			Code:      "ImgUpload",
			Type:      "MENU",
			ParentId:  intPointer(6),
			Path:      "/demo/upload",
			Icon:      "i-fe:image",
			Component: "/src/views/demo/upload/index.vue",
			KeepAlive: true,
			Show:      true,
			Enable:    true,
			SortOrder: 2,
		},
		{
			Name:      "个人资料",
			Code:      "UserProfile",
			Type:      "MENU",
			Path:      "/profile",
			Icon:      "i-fe:user",
			Component: "/src/views/profile/index.vue",
			Show:      false,
			Enable:    true,
			SortOrder: 99,
		},
		{
			Name:      "基础功能",
			Code:      "Base",
			Type:      "MENU",
			Path:      "/base",
			Icon:      "i-fe:grid",
			Show:      true,
			Enable:    true,
			SortOrder: 0,
		},
		{
			Name:      "基础组件",
			Code:      "BaseComponents",
			Type:      "MENU",
			ParentId:  intPointer(9),
			Path:      "/base/components",
			Icon:      "i-me:awesome",
			Component: "/src/views/base/index.vue",
			Show:      true,
			Enable:    true,
			SortOrder: 1,
		},
		{
			Name:      "Unocss",
			Code:      "Unocss",
			Type:      "MENU",
			ParentId:  intPointer(9),
			Path:      "/base/unocss",
			Icon:      "i-me:awesome",
			Component: "/src/views/base/unocss.vue",
			Show:      true,
			Enable:    true,
			SortOrder: 2,
		},
		{
			Name:      "KeepAlive",
			Code:      "KeepAlive",
			Type:      "MENU",
			ParentId:  intPointer(9),
			Path:      "/base/keep-alive",
			Icon:      "i-me:awesome",
			Component: "/src/views/base/keep-alive.vue",
			KeepAlive: true,
			Show:      true,
			Enable:    true,
			SortOrder: 3,
		},
		{
			Name:      "创建新用户",
			Code:      "AddUser",
			Type:      "BUTTON",
			ParentId:  intPointer(4),
			Show:      true,
			Enable:    true,
			SortOrder: 1,
		},
		{
			Name:      "图标 Icon",
			Code:      "Icon",
			Type:      "MENU",
			ParentId:  intPointer(9),
			Path:      "/base/icon",
			Icon:      "i-fe:feather",
			Component: "/src/views/base/unocss-icon.vue",
			Show:      true,
			Enable:    true,
			SortOrder: 0,
		},
		{
			Name:      "MeModal",
			Code:      "TestModal",
			Type:      "MENU",
			ParentId:  intPointer(9),
			Path:      "/testModal",
			Icon:      "i-me:dialog",
			Component: "/src/views/base/test-modal.vue",
			Show:      true,
			Enable:    true,
			SortOrder: 5,
		},
	}

	profiles := []model.Profile{
		{
			Gender:   0,
			Avatar:   "",
			UserId:   1,
			NickName: "Admin",
		},
	}

	roles := []model.Role{
		{
			Code:   "SUPER_ADMIN",
			Name:   "超级管理员",
			Enable: true,
		},
		{
			Code:   "USER",
			Name:   "普通用户",
			Enable: true,
		},
	}

	rolePermissions := []model.RolePermissionsPermission{
		{RoleId: 2, PermissionId: 1},
		{RoleId: 2, PermissionId: 2},
		{RoleId: 2, PermissionId: 3},
		{RoleId: 2, PermissionId: 4},
		{RoleId: 2, PermissionId: 5},
		{RoleId: 2, PermissionId: 9},
		{RoleId: 2, PermissionId: 10},
		{RoleId: 2, PermissionId: 11},
		{RoleId: 2, PermissionId: 12},
		{RoleId: 2, PermissionId: 14},
		{RoleId: 2, PermissionId: 15},
	}

	users := []model.User{
		{
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			Username:   "admin",
			Password:   "5f4dcc3b5aa765d61d8327deb882cf99", // 密码经过 MD5 加密
			Enable:     true,
		},
	}

	userRoles := []model.UserRolesRole{
		{UserId: 1, RoleId: 1},
		{UserId: 1, RoleId: 2},
	}

	db.Create(&permissions)
	db.Create(&profiles)
	db.Create(&roles)
	db.Create(&rolePermissions)
	db.Create(&users)
	db.Create(&userRoles)
}
