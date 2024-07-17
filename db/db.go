package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gin-naiveui/config"
	"gin-naiveui/model"

	"gorm.io/driver/postgres"
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

	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	fmt.Println(port)

	dsn := fmt.Sprintf(
		"host=naiva_postgres port=%d user=%s password=%s dbname=%s sslmode=disable",
		port,
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	)

	openDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   dbLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("db connection error is %s", err.Error())
	}

	openDb.AutoMigrate(&model.User{}, &model.Profile{}, &model.Role{}, &model.Permission{}, &model.UserRolesRole{}, &model.RolePermissionsPermission{})

	dbCon, err := openDb.DB()
	if err != nil {
		log.Fatalf("openDb.DB error is  %s", err.Error())
	}
	dbCon.SetMaxIdleConns(3)
	dbCon.SetMaxOpenConns(10)
	dbCon.SetConnMaxLifetime(time.Hour)
	Dao = openDb
}
