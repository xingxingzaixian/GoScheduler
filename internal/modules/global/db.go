package global

import "gorm.io/gorm"

var DB *gorm.DB

const (
	Page     = 1  // 当前页数
	PageSize = 20 // 每页多少条数据
)
