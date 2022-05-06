package model

import "gorm.io/gorm"

type Apply struct {
	gorm.Model
	ID       uint   `gorm:"primary_key" form:"id" json:"id"`
	UserId   uint   `form:"userId" json:"userId"`
	AffairId uint   `form:"affairsId" json:"affairsId"`
	PersonId uint   `form:"personId" json:"personId"`
	FileUrl  string `form:"fileUrl" json:"fileUrl"`
	Access   string `form:"access" json:"access"`
}

func (Apply) TableName() string {
	return "apply"
}
