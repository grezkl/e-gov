package main

import (
	"e-gov/database"
	_ "e-gov/docs"
	"e-gov/global"
	"e-gov/initialize"
	"e-gov/router"
	"e-gov/router/account"
	"e-gov/router/affair"
	"e-gov/router/apply"
	"e-gov/router/feedback"
	"e-gov/router/file"
	"e-gov/router/menu"
	"e-gov/router/news"
	"e-gov/router/person"
	"e-gov/router/role"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title e-gov server
// @version 1.0
// @description e-gov server

// @contact.name grezkl
// @contact.url http://github.com/grezkl
// @contact.email grezkl@protonmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9090
// @BasePath /

func main() {
	// 配置初始化
	initialize.InitConfig()

	// 数据库初始化
	database.Init()

	// 日志初始化，同时写入文件和控制台
	gin.DisableConsoleColor()
	date := time.Now().Format("20060102_150405")
	f, _ := os.Create(global.Settings.LogsAddress + "gin_" + date + ".log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	// 路由初始化
	r = router.Init(account.Router, role.Router, menu.Router, file.Router, affair.Router, person.Router, apply.Router, news.Router, feedback.Router)

	// 开启跨域中间件
	// r.Use(middleware.Cors())
	r.Use(cors.Default())

	// 启动
	err := r.Run(fmt.Sprintf("%s:%d", global.Settings.Host, global.Settings.Port))
	if err != nil {
		fmt.Println("服务器启动失败")
	}

}
