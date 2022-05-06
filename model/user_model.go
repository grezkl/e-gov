package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primary_key" form:"id" json:"id"`
	Username  string `gorm:"not null;unique" form:"username" json:"username" binding:"required"`
	Password  string `gorm:"default:'$2a$10$bPAPp3DbZuwjpHJshBqc9ubEbCylG4HNSQvnJbqRsYbCLHN7l8JJy'" form:"password" json:"password"`
	Nickname  string `form:"nickname" json:"nickname"`
	AvatarUrl string `form:"avatarUrl" json:"avatarUrl"`
	Email     string `form:"email" json:"email"`
	Phone     string `form:"phone" json:"phone"`
	Address   string `form:"address" json:"address"`
	Role      string `form:"role" json:"role" default:"ROLE_USER"`
	PersonId  uint   `form:"personId" json:"personId"`
}

type UserAuth struct {
	User  User   `form:"user" json:"user"`
	Token string `form:"token" json:"token"`
	Menus []Menu `form:"menus" json:"menus"`
	// Menus MenuList
}

func (User) TableName() string {
	return "user"
}

type UserLoginTemp struct {
	Username string `gorm:"not null;unique" form:"username" json:"username" binding:"required"`
	Password string `gorm:"not null" form:"password" json:"password" binding:"required"`
}

type UserRegisterTemp struct {
	Username string `gorm:"not null;unique" form:"username" json:"username" binding:"required"`
	Password string `gorm:"not null" form:"password" json:"password" binding:"required"`
	Nickname string `form:"nickname" json:"nickname"`
}
