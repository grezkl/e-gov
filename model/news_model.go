package model

import "gorm.io/gorm"

type News struct {
	gorm.Model
    ID          uint   `gorm:"primary_key" form:"id" json:"id"`
	Title       string `form:"title" json:"title"`
	Author      string `form:"author" json:"author"`
	Image       string `form:"image" json:"image"`
	Content     string `form:"content" json:"content"`
	Description string `form:"description" json:"description"`
}

func (News) TableName() string {
	return "news"
}
