package models

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	Name   string `json:"name" gorm:"type:varchar(64);not null"`
	Alias  string `json:"alias" gorm:"type:varchar(32);not null;default:''"`
	Port   int    `json:"port" gorm:"not null;default:5921"`
	Remark string `json:"remark" gorm:"type:varchar(100);not null;default:''"`
}

//func (h *Host) TableName() string {
//	return "tbl_host"
//}
