package model

type Person struct {
	ID       uint   `gorm:"primary_key" form:"id" json:"id"`
	Name     string `form:"name" json:"name" binding:"required"`
	Age      string `form:"age" json:"age" binding:"required"`
	Sex      string `form:"sex" json:"sex" binding:"required"`
	Birthday string `form:"birthday" json:"birthday" binding:"required"`
	Address  string `form:"address" json:"address" binding:"required"`
	Hometown string `form:"hometown" json:"hometown" binding:"required"`
	Identity string `form:"identity" json:"identity" binding:"required"`
	UserId   uint   `form:"userId" json:"userId" binding:"required"`
}

func (Person) TableName() string {
	return "person"
}
