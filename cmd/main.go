package main

import (
	"gin-naiveui/config"
	"gin-naiveui/db"
	"gin-naiveui/router"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var Loc, _ = time.LoadLocation("Asia/Shanghai")
	time.Local = Loc
	app := gin.Default()
	config.Init()
	db.Init()
	router.Init(app)
	app.Run(":3000")
}
