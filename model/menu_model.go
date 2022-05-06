package model

type Menu struct {
	ID          uint   `gorm:"primary_key" form:"id" json:"id"`
	Name        string `form:"name" json:"name"`
	Path        string `form:"path" json:"path"`
	Icon        string `form:"icon" json:"icon"`
	Description string `form:"description" json:"description"`
	Children    []Menu `gorm:"-" form:"children" json:"children"`
	Pid         string `form:"pid" json:"pid"`
	PagePath    string `form:"pagePath" json:"pagePath"`
}

func (Menu) TableName() string {
	return "menu"
}

// abandon
type MenuList struct {
	Items []Menu
}

func (menus *MenuList) Append(item Menu) []Menu {
	menus.Items = append(menus.Items, item)
	return menus.Items
}
