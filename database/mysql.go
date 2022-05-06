package database

import (
	"e-gov/global"
	"e-gov/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init() {
	var err error
	mysqlInfo := global.Settings.Mysqlinfo

	username := mysqlInfo.Username
	password := mysqlInfo.Password
	database := mysqlInfo.Database
	host := mysqlInfo.Host
	port := mysqlInfo.Port

	// loc 设置时区为本地，默认 UTC +0
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local",
		username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}

	global.DB = db
	global.DB.AutoMigrate(&model.User{}, &model.Role{}, &model.Menu{},
		&model.RoleMenu{}, &model.File{}, &model.Affair{}, &model.Feedback{},
		&model.Person{}, &model.News{}, &model.Apply{}, &model.Dict{})
}
