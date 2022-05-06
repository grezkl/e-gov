package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Affair struct {
	gorm.Model
	ID          uint            `gorm:"primary_key" form:"id" json:"id"`
	Name        string          `form:"name" json:"name"`
	Department  string          `form:"department" json:"department"`
	Theme       string          `form:"theme" json:"theme"`
	Cost        decimal.Decimal `form:"cost" json:"cost" sql:"type:decimal(20,8);"`
	State       bool            ` gorm:"default:true" form:"state" json:"state"`
	Description string          `form:"description" json:"description"`
	AuditId     uint            `form:"auditId" json:"auditId"`
}

func (Affair) TableName() string {
	return "affair"
}

type AffairRes struct {
	ID          uint   `form:"id" json:"id"`
	Name        string `form:"name" json:"name"`
	Department  string `form:"department" json:"department"`
	Theme       string `form:"theme" json:"theme"`
	Cost        int    `form:"cost" json:"cost"`
	State       bool   `form:"state" json:"state"`
	Description string `form:"description" json:"description"`
	AuditId     uint   `form:"auditId" json:"auditId"`
	Audit       string `form:"audit" json:"audit"`
}
