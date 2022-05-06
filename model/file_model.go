package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	ID     uint   `gorm:"primary_key" form:"id" json:"id"`
	Name   string `form:"name" json:"name"`
	Type   string `form:"type" json:"type"`
	Size   int64  `form:"size" json:"size"`
	Url    string `form:"url" json:"url"`
	MD5    string `form:"md5" json:"md5"`
	Enable bool   `form:"enable" json:"enable"`
	UserId uint   `form:"userId" json:"userId"`
}

func (File) TableName() string {
	return "file"
}
