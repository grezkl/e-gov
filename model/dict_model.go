package model

type Dict struct {
	Name  string `form:"name" json:"name"`
	Value string `form:"value" json:"value"`
	Type  string `form:"type" json:"type"`
}

func (Dict) TableName() string {
	return "dict"
}
