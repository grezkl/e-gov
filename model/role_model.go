package model

type Role struct {
	ID          uint   `gorm:"primary_key" form:"id" json:"id"`
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
	Flag        string `form:"flag" json:"flag"`
}

func (Role) TableName() string {
	return "role"
}
