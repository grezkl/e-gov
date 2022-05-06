package model

import "gorm.io/gorm"

type Feedback struct {
	gorm.Model
	ID           uint   `gorm:"primary_key" form:"id" json:"id"`
	Type         string `form:"type" json:"type"`
	Content      string `form:"content" json:"content"`
	UserId       uint   `form:"userId" json:"userId"`
	Responder    string `form:"responder" json:"responder"`
	ReplyContent string `form:"replyContent" json:"replyContent"`
	State        bool   ` gorm:"default:true" form:"state" json:"state"`
	Description  string `form:"description" json:"description"`
}

func (Feedback) TableName() string {
	return "feedback"
}
