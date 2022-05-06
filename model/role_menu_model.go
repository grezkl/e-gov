package model

type RoleMenu struct {
    RoleId uint `form:"roleId" json:"roleId"` 
    MenuId uint `form:"menuId" json:"menuId"`
}

func (RoleMenu) TableName() string {
	return "role_menu"
}
